package views

import (
	"net/http"

	"github.com/pchan37/studybuddy/src/lib/announcements"
	"github.com/pchan37/studybuddy/src/lib/security"
	"github.com/pchan37/studybuddy/src/lib/templateManager"
)

func RegisterPrivateViews() {
	http.HandleFunc("/", security.Authenticate(IndexPage))
	http.HandleFunc("/announcements", security.Authenticate(AnnouncementHandler))
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/page_not_found", http.StatusTemporaryRedirect)
	}
	templateManager.RenderTemplate(w, "index.tmpl", nil)
}

func AnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	session := security.GetSecuritySession(w, r)
	listOfAnnouncements, _ := announcements.GetAnnouncements()
	data := struct {
		Role          string
		Announcements []announcements.Announcement
	}{
		Role:          session.Values["role"].(string),
		Announcements: listOfAnnouncements,
	}
	templateManager.RenderTemplate(w, "announcement.tmpl", data)
}
