package views

import (
	"net/http"
	"time"

	"github.com/pchan37/studybuddy/src/lib/announcements"
)

func RegisterAPIViews() {
	http.HandleFunc("/save_announcements", SaveAnnouncements)
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
