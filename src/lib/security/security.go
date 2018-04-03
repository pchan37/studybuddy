package security

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var authKey, authKeyError = ioutil.ReadFile("../authKey")
var encryptKey, encryptKeyError = ioutil.ReadFile("../encryptKey")
var store = sessions.NewCookieStore(authKey, encryptKey)

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isLoggedIn(w, r) {
			h.ServeHTTP(w, r)
		}
		session, _ := store.Get(r, "security")
		session.Values["redirect-url"] = r.URL.Path
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	})
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) (loggedIn bool) {
	loggedIn = false
	session, _ := store.Get(r, "security")

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		loggedIn = true
	}
	return
}

func isSuccessful(success bool) bool {
	return success
}

func addFlashMessage(w http.ResponseWriter, r *http.Request, message string) {
	flashSession, _ := store.Get(r, "flashSession")
	flashSession.AddFlash(message)
	flashSession.Save(r, w)
}

func handleRedirect(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
	if redirectURL, ok := session.Values["redirect-url"].(string); ok && redirectURL != "" {
		session.Values["redirect-url"] = ""
		session.Save(r, w)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if isSuccessful(Login(&credential{r.FormValue("username"), r.FormValue("password"), ""})) {
		session, _ := store.Get(r, "security")
		session.Values["authenticated"] = true
		session.Save(r, w)
		addFlashMessage(w, r, "Successfully logged in!")
		handleRedirect(w, r, session)
	}
	addFlashMessage(w, r, "Incorrect username or password!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func validateRegisterCredentials(w http.ResponseWriter, r *http.Request, c *credential) {
	if IsRegistered(c.username) {
		addFlashMessage(w, r, "Username already taken!")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	} else if c.password != c.confirmationPassword {
		addFlashMessage(w, r, "Password does not match confirmation password!")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	userCredential := &credential{
		r.FormValue("username"),
		r.FormValue("password"),
		r.FormValue("confirmationPassword"),
	}
	log.Println("Validating...")
	validateRegisterCredentials(w, r, userCredential)
	log.Println("Successfully validated...")
	if isSuccessful(Register(userCredential)) {
		session, _ := store.Get(r, "security")
		session.Values["authenticated"] = true
		session.Save(r, w)
		addFlashMessage(w, r, "Successfully registered!")
		handleRedirect(w, r, session)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "security")

	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
