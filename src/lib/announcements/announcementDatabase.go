package announcements

import (
	"github.com/globalsign/mgo"
	"github.com/pchan37/studybuddy/src/lib/dbManager"
	"github.com/pchan37/studybuddy/src/utils"
)

type Announcement map[string]string

var announcementDB *dbManager.DBManager
var announcementCollection *mgo.Collection

func InitializeAnnouncementsDB(name string) *dbManager.DBManager {
	announcementDB = dbManager.New(name, "127.0.0.1:27017")
	announcementCollection = announcementDB.Database.C("announcements")
	return announcementDB
}

func GetAnnouncements() (result []Announcement, err error) {
	err = announcementDB.Database.C("announcements").Find(nil).All(&result)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return
}

func AddAnnouncements(announcement Announcement) {
	err := announcementDB.Database.C("announcements").Insert(announcement)
	utils.LogFatalOnError(err)
}
