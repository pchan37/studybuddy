package views

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/pchan37/studybuddy/src/lib/security"
	"github.com/pchan37/studybuddy/src/lib/templateManager"
)

var authKey, authKeyError = ioutil.ReadFile("../authKey")
var encryptKey, encryptKeyError = ioutil.ReadFile("../encryptKey")
var store = sessions.NewCookieStore(authKey, encryptKey)

func RegisterSecurityViews() {
	http.HandleFunc("/login", LoginDelegator)
	http.HandleFunc("/logout", security.LogoutHandler)
	http.HandleFunc("/register", RegisterDelegator)
}

func LoginDelegator(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if !security.IsLoggedIn(w, r) {
			flashSession, _ := store.Get(r, "flashSession")
			flashes := flashSession.Flashes()
			flashSession.Save(r, w)
			data := make(map[string]string)
			if len(flashes) > 0 {
				data["Messages"] = flashes[0].(string)
			} else {
				data["Messages"] = ""
			}
			templateManager.RenderTemplate(w, "login.tmpl", data)
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	case "POST":
		security.LoginHandler(w, r)
	}
}

func RegisterDelegator(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if !security.IsLoggedIn(w, r) {
			flashSession, _ := store.Get(r, "flashSession")
			flashes := flashSession.Flashes()
			flashSession.Save(r, w)
			data := make(map[string]string)
			if len(flashes) > 0 {
				data["Messages"] = flashes[0].(string)
			} else {
				data["Messages"] = ""
			}
			templateManager.RenderTemplate(w, "register.tmpl", data)
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	case "POST":
		security.RegisterHandler(w, r)
	}
}
