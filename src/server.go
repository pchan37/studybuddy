package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"

	"github.com/pchan37/studybuddy/src/lib/announcements"
	"github.com/pchan37/studybuddy/src/lib/dbManager"
	"github.com/pchan37/studybuddy/src/lib/security"
	"github.com/pchan37/studybuddy/src/lib/templateManager"
	"github.com/pchan37/studybuddy/src/utils"
	"github.com/pchan37/studybuddy/src/views"
)

type config struct {
	IncludePath string
	LayoutPath  string

	PrivateKeyPath string
	PublicKeyPath  string
}

func (c *config) logConfig() {
	log.Println("Include path:", c.IncludePath)
	log.Println("Layout path:", c.LayoutPath)
}

func loadConfig(filename string) (c *config) {
	file, err := os.Open(filename)
	utils.LogFatalfOnError("Unable to open config file: %s", err)
	decoder := json.NewDecoder(file)
	c = &config{}
	err = decoder.Decode(c)
	utils.LogFatalfOnError("Error unpacking config: %s", err)
	return
}

func main() {
	config := loadConfig("config.json")
	config.logConfig()
	templateManager.SetTemplateConfig(config.LayoutPath, config.IncludePath)
	templateManager.LoadTemplates()

	authManager := security.InitializeAuthManager("authentication")
	defer authManager.Close()

	announcementsDB := announcements.InitializeAnnouncementsDB("studybuddy")
	defer dbManager.Close(announcementsDB)

	server := http.Server{
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      context.ClearHandler(http.DefaultServeMux),
	}

	views.RegisterAPIViews()
	views.RegisterStaticViews()
	views.RegisterPrivateViews()
	views.RegisterPublicViews()
	views.RegisterSecurityViews()

	server.ListenAndServe()
}
