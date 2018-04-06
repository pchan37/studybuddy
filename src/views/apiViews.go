package views

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pchan37/studybuddy/src/lib/announcements"
	"github.com/pchan37/studybuddy/src/lib/security"
)

func RegisterAPIViews() {
	http.HandleFunc("/save_announcements", SaveAnnouncements)
	http.HandleFunc("/adjust_roles", AdjustRoles)
}

func SaveAnnouncements(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Redirect(w, r, "/page_not_found", http.StatusNotFound)
	case "POST":
		currentTime := time.Now()
		announcement := &announcements.Announcement{
			"Title": r.FormValue("title"),
			"Time":  "Posted on: " + currentTime.Format("Monday, January 2, 2006 3:04:05 PM MST"),
			"Body":  r.FormValue("body"),
		}
		announcements.AddAnnouncements(*announcement)
		http.Redirect(w, r, "/announcements", http.StatusTemporaryRedirect)
	}
}

func AdjustRoles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Redirect(w, r, "/page_not_found", http.StatusNotFound)
	case "POST":
		teacherAssistants := strings.Split(r.FormValue("teacher_assistants"), "\n")
		teachers := strings.Split(r.FormValue("teachers"), "\n")
		admins := strings.Split(r.FormValue("admins"), "\n")
		developers := strings.Split(r.FormValue("developers"), "\n")
		for _, teacherAssistant := range teacherAssistants {
			credential, ok := security.GetUserCredential(teacherAssistant)
			if ok {
				security.AddTeacherAssistant(credential)
			}
		}
		for _, teacher := range teachers {
			credential, ok := security.GetUserCredential(teacher)
			if ok {
				security.AddTeacher(credential)
			}
		}
		for _, admin := range admins {
			credential, ok := security.GetUserCredential(admin)
			if ok {
				security.AddAdmin(credential)
			}
		}
		for _, developer := range developers {
			credential, ok := security.GetUserCredential(developer)
			fmt.Println(credential)
			if ok {
				security.AddDeveloper(credential)
			}
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
