package models

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/sessions"
	"github.com/user/golang_signup/db"
	"github.com/user/golang_signup/tmpl"
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

func StartPage(w http.ResponseWriter, req *http.Request) {
	conditionsMap := map[string]interface{}{}
	session, err := store.Get(req, "session")
	fmt.Println(session)

	if err != nil {
		handleError(err, "Unable to retrieve session data!: %v", w)
		return
	}
	conditionsMap["username"] = session.Values["username"]
	if session != nil && conditionsMap["username"] == "" {
		conditionsMap["CheckSession"] = false
	} else if conditionsMap["username"] != ""{
		conditionsMap["CheckSession"] = true
	}

	if err := tmpl.StartPageTemplate.Execute(w, conditionsMap); err != nil {
		log.Println(err)
	}
}

func Dashboard(w http.ResponseWriter, req *http.Request) {
	conditionsMap := map[string]interface{}{}
	session, _ := store.Get(req, "session")

	if session != nil {
		conditionsMap["username"] = session.Values["username"]
		fmt.Println("conditionsMap -> username -> ", conditionsMap["username"])
		if conditionsMap["username"] == nil {
			http.Redirect(w, req, "/", http.StatusFound)
		}
	}

	log.Println(session)
	err := tmpl.DashboardTemplate.Execute(w, conditionsMap)
	if err != nil {
		log.Println(err)
	}

}

func LoginUser(w http.ResponseWriter, req *http.Request) {
	conditionsMap := map[string]interface{}{}

	// check if session is active
	session, _ := store.Get(req, "session")
	if session != nil {
		conditionsMap["username"] = session.Values["username"]
		conditionsMap["password"] = session.Values["password"]
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	login := req.FormValue("login")

	log.Println("Login in:", username)
	log.Println("Login password:", password)

	if login != "" {
		if username == "" || password == "" {
			conditionsMap["LoginError"] = true
		} else {
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

func LogoutUser(w http.ResponseWriter, req *http.Request) {
	// Read from session
	session, _ := store.Get(req, "session")

	// Remove the username
	session.Values["username"] = ""
	err := session.Save(req, w)

	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, req, "/", http.StatusFound)
}

func SignUpUser(w http.ResponseWriter, req *http.Request) {
	conditionsMap := map[string]interface{}{}

	// check if session is active
	session, _ := store.Get(req, "session")
	if session != nil {
		conditionsMap["username"] = session.Values["username"]
		conditionsMap["password"] = session.Values["password"]
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	signup := req.FormValue("signup")

	if signup != "" {
		if username == "" || password == "" {
			conditionsMap["SignUpError"] = true
		} else {
			conditionsMap["SignUpError"] = false

			//ДОПИСАТЬ ПРОВЕРКУ НА КОРРЕКТНОСТЬ ВВОДА ДАННЫХ
			if err := checkUsername(username); err != true {

			}
			if err := checkPassword(password); err != true {

			}

			var user db.User
			user.Id = bson.NewObjectId()
			user.Username = username
			user.Password = password

			if err := db.CreateUser(user); err != nil {
				fmt.Println("Error with insert")
				conditionsMap["SignUpSuccessful"] = false
				return
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
	err := tmpl.SignUpTemplate.Execute(w, conditionsMap)
	if err != nil {
		log.Println(err)
	}
}

func checkPassword(password string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", password); !ok {
		return false
	}
	return true
}

func checkUsername(username string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", username); !ok {
		return false
	}
	return true
}
