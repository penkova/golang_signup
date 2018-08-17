package main

import (
	"net/http"
	"os"

	"log"

	"github.com/gorilla/context"
	"github.com/user/golang_signup/service/config"
	"github.com/user/golang_signup/service/models"
)

var (
	f    *os.File
	port = ":8080"
)

func init() {
	logFile := config.ReadLogFileConfig()
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	f = file
	//set output of logs to f
	log.SetOutput(f)
}

func main() {
	defer end()
	log.Printf("Listening to http://localhost%s", port)
	// Here we are instantiating the gorilla/mux router
	//r := mux.NewRouter()

	//r.HandleFunc("/", models.StartPage).Methods("POST")
	http.HandleFunc("/", models.StartPage)
	http.HandleFunc("/signup", models.SignUpUser)
	http.HandleFunc("/login", models.LoginUser)
	http.HandleFunc("/dashboard", models.Dashboard)
	http.HandleFunc("/logout", models.LogoutUser)

	// Our application will run on port 3030. Here we declare the port and pass in our router.
	http.ListenAndServe(port, context.ClearHandler(http.DefaultServeMux))
	//log.Println(http.ListenAndServe(port, nil))
}

func end() {
	if f != nil {
		defer f.Close()
	}
}
