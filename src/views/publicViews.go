package views

import (
	"net/http"

	"github.com/pchan37/studybuddy/src/lib/security"
	"github.com/pchan37/studybuddy/src/lib/templateManager"
)

func RegisterPublicViews() {
	http.HandleFunc("/page_not_found", NotFoundPage)
	http.HandleFunc("/", security.Authenticate(IndexPage))
}

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	templateManager.RenderTemplate(w, "404.tmpl", nil)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/page_not_found", http.StatusTemporaryRedirect)
	}
	data := map[string]bool{"IsLoggedIn": security.IsLoggedIn(w, r)}
	templateManager.RenderTemplate(w, "index.tmpl", data)
}
