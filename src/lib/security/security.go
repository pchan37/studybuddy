package security

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/pchan37/studybuddy/src/lib/templateManager"
)

var authKey, authKeyError = ioutil.ReadFile("../authKey")
var encryptKey, encryptKeyError = ioutil.ReadFile("../encryptKey")
var store = sessions.NewCookieStore(authKey, encryptKey)

func GetSecuritySession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "security")
	return session
}

func Authenticate(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsLoggedIn(w, r) {
			h(w, r)
		}
		session, _ := store.Get(r, "security")
		session.Values["redirect-url"] = r.URL.Path
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	})
}

func Authorize(h http.HandlerFunc, isAuthorized func(string) bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "security")
		role := session.Values["role"].(string)
		if isAuthorized(role) {
			h(w, r)
		}
		http.Redirect(w, r, "/permission_denied", http.StatusTemporaryRedirect)
	})
}

func IsLoggedIn(w http.ResponseWriter, r *http.Request) (loggedIn bool) {
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
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	userCredential := &credential{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	if isSuccessful(Login(userCredential)) {
		fullCredential, _ := GetUserCredential(userCredential.Username)
		session, _ := store.Get(r, "security")
		session.Values["authenticated"] = true
		session.Values["role"] = fullCredential.Role
		session.Save(r, w)
		handleRedirect(w, r, session)
	}
	addFlashMessage(w, r, "Incorrect username or password!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func validateRegisterCredentials(w http.ResponseWriter, r *http.Request, c *credential) {
	if IsRegistered(c.Username) {
		addFlashMessage(w, r, "Username already taken!")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	} else if c.Password != c.ConfirmationPassword {
		addFlashMessage(w, r, "Password does not match confirmation password!")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	userCredential := &credential{
		Username:             r.FormValue("username"),
		Password:             r.FormValue("password"),
		ConfirmationPassword: r.FormValue("confirmationPassword"),
	}
	validateRegisterCredentials(w, r, userCredential)
	if isSuccessful(Register(userCredential)) {
		fullCredential, _ := GetUserCredential(userCredential.Username)
		session, _ := store.Get(r, "security")
		session.Values["authenticated"] = true
		session.Values["role"] = fullCredential.Role
		session.Save(r, w)
		handleRedirect(w, r, session)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "security")

	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func NotAuthorizedHandler(w http.ResponseWriter, r *http.Request) {
	templateManager.RenderTemplate(w, "403.tmpl", nil)
}
