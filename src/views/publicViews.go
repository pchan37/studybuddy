package views

import (
	"net/http"

	"github.com/pchan37/studybuddy/src/lib/security"
	"github.com/pchan37/studybuddy/src/lib/templateManager"
)

func RegisterPublicViews() {
	http.HandleFunc("/page_not_found", NotFoundPage)
	http.HandleFunc("/permission_denied", security.NotAuthorizedHandler)
	http.HandleFunc("/teacher_only", security.Authorize(security.Authenticate(IndexPage), security.IsTeacher))
}

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	templateManager.RenderTemplate(w, "404.tmpl", nil)
}
