package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// User represents database entity.
type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"password" bson:"password"`
}

var db *mgo.Database

// Establish a connection to MongoDB database.
func init() {
	session, err := mgo.Dial("localhost:27018")

	if err != nil {
		fmt.Println("Error")
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db = session.DB("authmongo")
	fmt.Println("Woow")
}

// CollectionUsers - connection to usersdb collections
func CollectionUsers() *mgo.Collection {
	return db.C("usersdb")
}

// CreateUser - create User
func CreateUser(user User) error {
	return CollectionUsers().Insert(user)
}

// FindUser - finding user by username and password
func FindUser(username string, password string) (*User, error) {
	res := User{}
	err := CollectionUsers().Find(bson.M{
		"username": username,
		"password": password,
	}).One(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
