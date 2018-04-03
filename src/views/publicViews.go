package views

import (
	"net/http"

	"github.com/pchan37/studybuddy/src/lib/security"
	"github.com/pchan37/studybuddy/src/lib/templateManager"
)

func RegisterPublicViews() {
	http.HandleFunc("/", security.Authenticate(IndexPage))
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	templateManager.RenderTemplate(w, "index.tmpl", nil)
}
