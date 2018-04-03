package security

import (
	"crypto/sha256"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/pchan37/studybuddy/src/utils"
)

var collection *mgo.Collection

type credential struct {
	username             string
	password             string
	confirmationPassword string
	role                 string
}

func InitializeAuthManager(databaseName string) (session *mgo.Session) {
	host := []string{"127.0.0.1:27017"}
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: host,
	})
	utils.PanicOnError(err)
	log.Print("Creating database...")
	database := session.DB(databaseName)
	collection = database.C("authentication")
	return
}

func isFatalError(err error) bool {
	return err != nil && err.Error() != "not found"
}

func getHashedPassword(password string) string {
	preHash := sha256.Sum256([]byte(password))
	passwordHashBytes, err := bcrypt.GenerateFromPassword(preHash[:], 12)
	utils.PanicOnError(err)
	return string(passwordHashBytes)
}

func IsRegistered(username string) (registered bool) {
	query := bson.M{"username": username}
	count, _ := collection.Find(query).Count()
	registered = count != 0
	return
}

func Register(c *credential) (success bool) {
	if IsRegistered(c.username) || c.password != c.confirmationPassword {
		success = false
	} else {
		hashedPassword := getHashedPassword(c.password)
		collection.Insert(bson.M{"username": c.username, "password": hashedPassword, "role": "student"})
		success = true
	}
	return
}

func Login(c *credential) (success bool) {
	hashedPassword := getHashedPassword(c.password)
	query := bson.M{"username": c.username, "password": hashedPassword}
	credentialFound := credential{}
	if !IsRegistered(c.username) {
		success = false
	} else if err := collection.Find(query).One(&credentialFound); !isFatalError(err) {
		success = true
	}
	return
}

func IsStudent(c *credential) bool {
	return c.role == "student"
}

func IsTeacherAssistant(c *credential) bool {
	return c.role == "teacher_assistant"
}

func IsTeacher(c *credential) bool {
	return c.role == "teacher"
}

func IsDeveloper(c *credential) bool {
	return c.role == "developer"
}

func AddStudent(c *credential) {
	c.role = "student"
}

func DropStudent(c *credential) {
	c.role = ""
}

func AddTeacherAssistant(c *credential) {
	c.role = "teacher_assistant"
}

func DropTeacherAssistant(c *credential) {
	c.role = ""
}

func AddTeacher(c *credential) {
	c.role = "teacher"
}

func DropTeacher(c *credential) {
	c.role = ""
}
