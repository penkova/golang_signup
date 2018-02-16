package main

import (
	"net/http"

	"github.com/user/golang_signup/models"

	//"github.com/gorilla/mux"
	"github.com/gorilla/context"
)

const(
	port = ":8080"
)

func main() {
	// Here we are instantiating the gorilla/mux router
	//r := mux.NewRouter()

	//r.HandleFunc("/", models.StartPage).Methods("POST")
	http.HandleFunc("/", models.StartPage)
	http.HandleFunc("/signup", models.SignUpUser)
	http.HandleFunc("/login", models.LoginUser)
	http.HandleFunc("/logout", models.LogoutUser)

	// Our application will run on port 3030. Here we declare the port and pass in our router.
	http.ListenAndServe(port, context.ClearHandler(http.DefaultServeMux))
}











