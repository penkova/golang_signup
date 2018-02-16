package models

import (
	"fmt"
	"net/http"
	//"github.com/user/golang_signup/db"
	"regexp"
	"github.com/gorilla/sessions"

	"io"
)

func handleError(err error, message string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(message, err)))
}

var store = sessions.NewCookieStore([]byte("secret-password"))

func StartPage(w http.ResponseWriter, req *http.Request){
	session, _:= store.Get(req, "session")
	if req.FormValue("username") != ""{
		session.Values["username"] = req.FormValue("username")
	}
	session.Save(req, w)
	io.WriteString(w, `<!DOCTYPE html>
	<html>
		<head>
			<title>Home Page</title>
		</head>
		<body>
			<div style="display: flex; flex-direction: column; align-items: center" >
				<h1>LOGIN</h1>
				<a href="/login">Login</a>
				<a href="/signup">Sign Up</a>
			</div>
		</body>
	</html>
	`)
}

func SignUpUser(w http.ResponseWriter, req *http.Request){

}

func LoginUser(w http.ResponseWriter, req *http.Request){
	session, _:= store.Get(req, "session")
	if req.FormValue("username") != ""{
		session.Values["username"] = req.FormValue("username")
	}
	if req.FormValue("password") != ""{
		session.Values["password"] = req.FormValue("password")
	}
	session.Save(req, w)
	io.WriteString(w, `<!DOCTYPE html>
	<html>
		<head>
			<title>Home Page</title>
		</head>
		<body>
			<div style="display: flex; flex-direction: column; align-items: center" >
				<h1>SIGN IN</h1>
				<form method="POST">
						`+fmt.Sprint(session.Values["username"])+ `
					<label for="username">Username:</label>
					<input id="username" name="username" value="">
						`+fmt.Sprint(session.Values["password"])+ `
					<label for="pass1">Password:</label>
					<input type="password" id="pass1" name="password" value="">
					<label for="pass1">Password:</label>
					<input type="password" id="pass2" name="password" value="">
					<input type="submit">
				</form>
			</div>
		</body>
	</html>`)
}

func LogoutUser(w http.ResponseWriter, req *http.Request){

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