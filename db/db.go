package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"password" bson:"password"`
}

var db *mgo.Database

func init() {
	session, err := mgo.Dial("localhost:27018")

	if err != nil {
		fmt.Println("Error")
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db = session.DB("authmongo")
	fmt.Println("Woow")
}

func CollectionUsers() *mgo.Collection {
	return db.C("usersdb")
}

func CreateUser(user User) error {
	return CollectionUsers().Insert(user)
}
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