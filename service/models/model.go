package models

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/sessions"
	"github.com/user/golang_signup/service/db"
	"github.com/user/golang_signup/service/tmpl"
	"gopkg.in/mgo.v2/bson"
)

var encryptionKey = "something-secret"
var store = sessions.NewCookieStore([]byte(encryptionKey))

func init() {
	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   60 * 2, // 2 min. After 2 minutes of authentication, the session will end
		HttpOnly: true,
	}
}

func handleError(err error, message string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(message, err)))
}

// StartPage start page
func StartPage(w http.ResponseWriter, req *http.Request) {
	conditionsMap := map[string]interface{}{}
	session, err := store.Get(req, "session")
	fmt.Println(session)

	if err != nil {
		handleError(err, "Unable to retrieve session data!: %v", w)
		return
	}
	conditionsMap["username"] = session.Values["username"]
	if session != nil && conditionsMap["username"] == nil {
		fmt.Println("nil")
		if session.Values["username"] == "" {
			conditionsMap["CheckSession"] = false
		} else {
			conditionsMap["CheckSession"] = true
		}
		conditionsMap["CheckSession"] = false
	} else if conditionsMap["username"] != nil {
		fmt.Println("no nil")
		conditionsMap["CheckSession"] = true
		if session.Values["username"] == "" {
			conditionsMap["CheckSession"] = false
		} else {
			conditionsMap["CheckSession"] = true
		}
	}

	if err := tmpl.StartPageTemplate.Execute(w, conditionsMap); err != nil {
		log.Println(err)
	}
}

// Dashboard - user page
func Dashboard(w http.ResponseWriter, req *http.Request) {
	conditionsMap := map[string]interface{}{}
	session, _ := store.Get(req, "session")
	// check if session is active. If session is inative - redirect to start page
	if session != nil {
		conditionsMap["username"] = session.Values["username"]
		if conditionsMap["username"] == nil {
			http.Redirect(w, req, "/", http.StatusFound)
		}
	}

	err := tmpl.DashboardTemplate.Execute(w, conditionsMap)
	if err != nil {
		log.Println(err)
	}

}

// LoginUser - authentication user
func LoginUser(w http.ResponseWriter, req *http.Request) {
	conditionsMap := map[string]interface{}{}

	// check if session is active
	session, _ := store.Get(req, "session")
	if session != nil {
		conditionsMap["username"] = session.Values["username"]
		conditionsMap["password"] = session.Values["password"]
	}
	// Obtaining values
	username := req.FormValue("username")
	password := req.FormValue("password")
	login := req.FormValue("login")

	log.Println("Login in:", username)
	log.Println("Login password:", password)
	// If you click on "let me in!"
	if login != "" {
		// If you do not fill one of the fields
		if username == "" || password == "" {
			conditionsMap["LoginError"] = true
		} else {
			// Check for this user in the database
			if _, err := db.FindUser(username, password); err != nil {
				fmt.Println("Error. Not USER")
				conditionsMap["UserNotDatabase"] = true
			} else {
				conditionsMap["UserNotDatabase"] = false

				// Create a new session and redirect to dashboard
				session, _ := store.New(req, "session")
				fmt.Println("My username %s" + username)

				session.Values["username"] = username
				err := session.Save(req, w)
				if err != nil {
					log.Println(err)
				}
				http.Redirect(w, req, "/dashboard", http.StatusFound)
			}
			conditionsMap["LoginError"] = false
		}
	}

	err := tmpl.LoginUserTemplate.Execute(w, conditionsMap)
	if err != nil {
		log.Println(err)
	}
}

// LogoutUser Session => nil
func LogoutUser(w http.ResponseWriter, req *http.Request) {
	// Read from session
	session, _ := store.Get(req, "session")
	// Remove the username
	session.Values["username"] = ""
	err := session.Save(req, w)

	if err != nil {
		log.Println(err)
	}
	// Redirect to start page
	http.Redirect(w, req, "/", http.StatusFound)
}

// SignUpUser - user registration
func SignUpUser(w http.ResponseWriter, req *http.Request) {
	conditionsMap := map[string]interface{}{}

	// check if session is active
	session, _ := store.Get(req, "session")
	if session != nil {
		conditionsMap["username"] = session.Values["username"]
		conditionsMap["password"] = session.Values["password"]
	}
	// Obtaining values
	username := req.FormValue("username")
	password := req.FormValue("password")
	signup := req.FormValue("signup")
	// If you click on "sign up"
	if signup != "" {
		// If you do not fill one of the fields
		if username == "" || password == "" {
			conditionsMap["SignUpError"] = true
		} else {
			conditionsMap["SignUpError"] = false

			// Check for correctness of data entry
			checkName := checkUsername(username)
			checkPass := checkPassword(password)
			if !checkName || !checkPass {
				conditionsMap["ErrorCheckUser"] = true
			} else {
				conditionsMap["ErrorCheckUser"] = false
				// Create new user
				var user db.User
				user.ID = bson.NewObjectId()
				user.Username = username
				user.Password = password
				// Creating new user in db
				if err := db.CreateUser(user); err != nil {
					fmt.Println("Error with insert")
					conditionsMap["SignUpSuccessful"] = false
				} else {
					conditionsMap["SignUpSuccessful"] = true
				}

				log.Println("Sign in username:", username)
				log.Println("Sign in password:", password)

				// Create a new session and redirect to login
				session, _ := store.New(req, "session")
				fmt.Println("My username %s" + username)

				session.Values["username"] = ""
				err := session.Save(req, w)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	err := tmpl.SignUpTemplate.Execute(w, conditionsMap)
	if err != nil {
		log.Println(err)
	}
}

// checkPassword - check password for the contents of letters and numbers (from 4 to 16 characters)
func checkPassword(password string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", password); !ok {
		return false
	}
	return true
}

// checkUsername - check username for the contents of letters and numbers (from 4 to 16 characters)
func checkUsername(username string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", username); !ok {
		return false
	}
	return true
}
