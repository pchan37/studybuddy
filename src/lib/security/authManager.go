package security

import (
	"crypto/sha256"
	"fmt"
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

func GetUserCredential(username string) (*credential, bool) {
	query := bson.M{"username": username}
	userCredential := credential{}
	err := collection.Find(query).One(&userCredential)
	if err == nil {
		return &userCredential, true
	}
	return nil, false
}

func IsRegistered(username string) (registered bool) {
	_, registered = GetUserCredential(username)
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
	if !IsRegistered(c.Username) {
		success = false
	} else if userCredential, ok := GetUserCredential(c.Username); ok {
		preHash := getPreHashedPassword(c.Password)
		err := bcrypt.CompareHashAndPassword([]byte(userCredential.Password), preHash[:])
		success = err == nil
	}
	return
}

func IsStudent(username string) bool {
	if userCredential, success := GetUserCredential(username); success {
		fmt.Println(userCredential.Role)
		return userCredential.Role == "student"
	}
	return false
}

func IsTeacherAssistant(username string) bool {
	if userCredential, success := GetUserCredential(username); success {
		return userCredential.Role == "teacher_assistant"
	}
	return false
}

func IsTeacher(username string) bool {
	if userCredential, success := GetUserCredential(username); success {
		return userCredential.Role == "teacher"
	}
	return false
}

func IsAdmin(username string) bool {
	if userCredential, success := GetUserCredential(username); success {
		return userCredential.Role == "admin"
	}
	return false
}

func IsDeveloper(username string) bool {
	if userCredential, success := GetUserCredential(username); success {
		return userCredential.Role == "developer"
	}
	return false
}

func AddStudent(c *credential) bool {
	if userCredential, success := GetUserCredential(c.Username); success {
		selector := bson.M{"username": userCredential.Username}
		updator := bson.M{"$set": bson.M{"role": "student"}}
		collection.Update(selector, updator)
		return true
	}
	return false
}

func DropStudent(c *credential) bool {
	if userCredential, success := GetUserCredential(c.Username); success {
		selector := bson.M{"username": userCredential.Username}
		updator := bson.M{"$set": bson.M{"role": ""}}
		collection.Update(selector, updator)
		return true
	}
	return false
}

func AddTeacherAssistant(c *credential) bool {
	if userCredential, success := GetUserCredential(c.Username); success {
		selector := bson.M{"username": userCredential.Username}
		updator := bson.M{"$set": bson.M{"role": "teacher_assistant"}}
		collection.Update(selector, updator)
		return true
	}
	return false
}

func DropTeacherAssistant(c *credential) bool {
	if userCredential, success := GetUserCredential(c.Username); success {
		selector := bson.M{"username": userCredential.Username}
		updator := bson.M{"$set": bson.M{"role": ""}}
		collection.Update(selector, updator)
		return true
	}
	return false
}

func AddTeacher(c *credential) bool {
	if userCredential, success := GetUserCredential(c.Username); success {
		selector := bson.M{"username": userCredential.Username}
		updator := bson.M{"$set": bson.M{"role": "teacher"}}
		collection.Update(selector, updator)
		return true
	}
	return false
}

func DropTeacher(c *credential) bool {
	if userCredential, success := GetUserCredential(c.Username); success {
		selector := bson.M{"username": userCredential.Username}
		updator := bson.M{"$set": bson.M{"role": ""}}
		collection.Update(selector, updator)
		return true
	}
	return false
}

func AddAdmin(c *credential) bool {
	if userCredential, success := GetUserCredential(c.Username); success {
		selector := bson.M{"username": userCredential.Username}
		updator := bson.M{"$set": bson.M{"role": "admin"}}
		collection.Update(selector, updator)
		return true
	}
	return false
}

func DropAdmin(c *credential) bool {
	if userCredential, success := GetUserCredential(c.Username); success {
		selector := bson.M{"username": userCredential.Username}
		updator := bson.M{"$set": bson.M{"role": ""}}
		collection.Update(selector, updator)
		return true
	}
	return false
}

func AddDeveloper(c *credential) bool {
	if userCredential, success := GetUserCredential(c.Username); success {
		selector := bson.M{"username": userCredential.Username}
		updator := bson.M{"$set": bson.M{"role": "developer"}}
		collection.Update(selector, updator)
		return true
	}
	return false
}

func DropDeveloper(c *credential) bool {
	if userCredential, success := GetUserCredential(c.Username); success {
		selector := bson.M{"username": userCredential.Username}
		updator := bson.M{"$set": bson.M{"role": ""}}
		collection.Update(selector, updator)
		return true
	}
	return false
}
