package models

import (
	"fmt"
	"net/http"
	"regexp"
	"log"

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
		MaxAge:   3600 * 3, // 3 hours
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
	if session != nil {
		conditionsMap["username"] = session.Values["username"]
		fmt.Println(conditionsMap["username"])



		// Написать проверку на то, что ты уже вышел из профиля... что-то с сессиями
		if conditionsMap["username"]  != nil {
			conditionsMap["CheckSession"] = true
		} else{
			conditionsMap["CheckSession"] = false
		}
	}
	if err := tmpl.StartPageTemplate.Execute(w, conditionsMap); err != nil {
		log.Println(err)
	}
}

func SignUpUser(w http.ResponseWriter, req *http.Request) {
	conditionsMap := map[string]interface{}{}

	// check if session is active
	session, _ := store.Get(req, "session")
	if session != nil {
		conditionsMap["username"] = session.Values["username"]
		conditionsMap["password"]= session.Values["password"]
	}



//НАПИСАТЬ ПРОВЕРКУ НА КОРРЕКТНОСТЬ ВВОДА ДАННЫХ



	username := req.FormValue("username")
	password := req.FormValue("password")
	signup := req.FormValue("signup")

	fmt.Println([]byte(password))

	if signup != "" {
		if username == "" || password == ""{
			conditionsMap["SignUpError"] = true
		} else {
			conditionsMap["SignUpError"] = false

			var user db.User
			user.Id = bson.NewObjectId()
			user.Username = username
			user.Password = password

			if err := db.CreateUser(user); err != nil {
				fmt.Println("Error with insert")
				handleError(err, "Failed to load database items: %v", w)
				conditionsMap["SignUpSuccessful"] = false
				return
			} else{
				conditionsMap["SignUpSuccessful"] = true
			}

			log.Println("Sign in username:", username)
			log.Println("Sign in password:", password)

			// Create a new session and redirect to login
			session, _ := store.New(req, "session")
			fmt.Println("My username %s"+ username)

			session.Values["username"] = username
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

func Dashboard(w http.ResponseWriter, req *http.Request) {
	conditionsMap := map[string]interface{}{}
	session, _ := store.Get(req, "session")
// ДОПИСАТЬ ЧТО БЫ ПРИ ЛОГАУТ ПОКАЗЫВАЛА СТРАНИЦА ЧТО-ТО ДРУГО И ПРИДУМАТЬ ЧТО
	if session != nil {
		conditionsMap["username"] = session.Values["username"]
		//conditionsMap["body"] = ""
	}
	err := tmpl.DashboardTemplate.Execute(w, conditionsMap)
	if err != nil {
		log.Println(err)
	}

}

func LoginUser(w http.ResponseWriter, req *http.Request){
	conditionsMap := map[string]interface{}{}

	// check if session is active
	session, _ := store.Get(req, "session")
	if session != nil {
		conditionsMap["username"] = session.Values["username"]
		conditionsMap["password"]= session.Values["password"]
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	login := req.FormValue("login")

	fmt.Println([]byte(password))

	if login != "" {
		if username == "" || password == ""{
			conditionsMap["LoginError"] = true
		} else {
			conditionsMap["SignUpError"] = false

			_, err := db.FindUser(username, password)
			if err != nil {
				fmt.Println("Error. not user")
				handleError(err, "Failed to load database items: %v", w)



				//ДОПИСАТЬ ЧТО БУДЕТ ЕСЛИ НЕТ ФАЙЛА



			return
			}

			log.Println("Login in:", username)
			log.Println("Login password:", password)

			// Create a new session and redirect to dashboard
			session, _ := store.New(req, "session")
			fmt.Println("My username %s"+ username)

			session.Values["username"] = username
			error := session.Save(req, w)
			if error != nil {
				log.Println(err)
			}
			http.Redirect(w, req, "/dashboard", http.StatusFound)
		}
	}


	err := tmpl.LoginUserTemplate.Execute(w, conditionsMap)
	if err != nil {
		log.Println(err)
	}
}
func LogoutUser(w http.ResponseWriter, req *http.Request){
	//read from session
	session, _ := store.Get(req, "session")

	// remove the username
	session.Values["username"] = ""
	err := session.Save(req, w)

	if err != nil {
		log.Println(err)
	}

	w.Write([]byte("Logged out!"))
}

//func SignUpUser(w http.ResponseWriter, req *http.Request) {
	//
	//	conditionsMap := map[string]interface{}{}
	//
	//	// check if session is active
	//	session, _ := store.Get(req, "session")
	//	if session != nil {
	//		conditionsMap["username"] = session.Values["username"]
	//		conditionsMap["password"]= session.Values["password"]
	//	}
	//
	//	username := req.FormValue("username")
	//	password := req.FormValue("password")
	//	signup := req.FormValue("signup")
	//
	//	fmt.Println([]byte(password))
	//
	//	if signup != "" {
	//		if username == "" || password == ""{
	//			conditionsMap["SignUpError"] = true
	//		} else {
	//			conditionsMap["SignUpError"] = false
	//
	//			var user db.User
	//			user.Id = bson.NewObjectId()
	//			user.Username = username
	//			user.Password = password
	//
	//			if err := db.CreateUser(user); err != nil {
	//				fmt.Println("Error with insert")
	//				handleError(err, "Failed to load database items: %v", w)
	//				return
	//			}
	//
	//
	//			hashedPasswordFromDatabase := []byte("$2a$10$4Yhs5bfGgp4vz7j6ScujKuhpRTA4l4OWg7oSukRbyRN7dc.C1pamu")
	//			if err := bcrypt.CompareHashAndPassword(hashedPasswordFromDatabase, []byte(password)); err != nil {
	//				log.Println("Either username or password is wrong")
	//
	//			} else {
	//				log.Println("Sign in :", username)
	//				log.Println("Sign in :", password)
	//				conditionsMap["Username"] = username
	//
	//				// Create a new session and redirect to dashboard
	//				session, _ := store.New(req, "session")
	//				fmt.Println("sdfssdf %s"+ username)
	//
	//				session.Values["username"] = username
	//				err := session.Save(req, w)
	//
	//				if err != nil {
	//					log.Println(err)
	//				}
	//
	//				http.Redirect(w, req, "/startpage", http.StatusFound)
	//			}
	//		}
	//	}
	//
	//
	//	err := tmpl.SignUpTemplate.Execute(w, conditionsMap)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//
	//}


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
