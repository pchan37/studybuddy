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
	Username             string
	Password             string
	ConfirmationPassword string
	Role                 string
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

func getPreHashedPassword(password string) [32]byte {
	return sha256.Sum256([]byte(password))
}

func getHashedPassword(password string) string {
	preHash := getPreHashedPassword(password)
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
	if IsRegistered(c.Username) || c.Password != c.ConfirmationPassword {
		success = false
	} else {
		hashedPassword := getHashedPassword(c.Password)
		collection.Insert(bson.M{"username": c.Username, "password": hashedPassword, "role": "student"})
		success = true
	}
	return
}

func Login(c *credential) (success bool) {
	query := bson.M{"username": c.Username}
	credentialFound := credential{}
	if !IsRegistered(c.Username) {
		success = false
	} else if err := collection.Find(query).One(&credentialFound); !isFatalError(err) {
		preHash := getPreHashedPassword(c.Password)
		err = bcrypt.CompareHashAndPassword([]byte(credentialFound.Password), preHash[:])
		success = err == nil
	}
	return
}

func IsStudent(c *credential) bool {
	return c.Role == "student"
}

func IsTeacherAssistant(c *credential) bool {
	return c.Role == "teacher_assistant"
}

func IsTeacher(c *credential) bool {
	return c.Role == "teacher"
}

func IsDeveloper(c *credential) bool {
	return c.Role == "developer"
}

func AddStudent(c *credential) {
	c.Role = "student"
}

func DropStudent(c *credential) {
	c.Role = ""
}

func AddTeacherAssistant(c *credential) {
	c.Role = "teacher_assistant"
}

func DropTeacherAssistant(c *credential) {
	c.Role = ""
}

func AddTeacher(c *credential) {
	c.Role = "teacher"
}

func DropTeacher(c *credential) {
	c.Role = ""
}
