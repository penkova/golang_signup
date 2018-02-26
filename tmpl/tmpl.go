package tmpl

import (
	"html/template"
)

const StartPage = `<!DOCTYPE html>
	<html>
		<head>
			<title>Home Page</title>
		</head>
		<body>
			<div style="display: flex; flex-direction: column; align-items: center" >
		{{if .CheckSession}}	
			<p style="color:blue ; text-transform: uppercase">you are already authenticated. Go to your <a href="/dashboard">profile!</a></p>
		{{else}}
				<div style="display: flex; flex-direction: column; align-items: center" >
				<h1>LOGIN</h1>
				<a href="/login">Login</a>
				<a href="/signup">Sign Up</a>
			</div>
		{{end}}	
		</div>
		</body>
	</html>`

const SingUpPage = `<!DOCTYPE html>
	<html>
	<head>
		<title>SIGN IN</title>
	</head>
	<body>
		<div style="display: flex; flex-direction: column; align-items: center" >
			<h1>SIGN IN</h1>
			{{if .SignUpError}}	<p style="color:red ; text-transform: uppercase">No username or password! Fill in required fields </p>{{end}}
			{{if .ErrorCheckUser}}	<p style="color:red ; text-decoration: underline; text-transform: uppercase">Username and password can only contain letters AND numbers</p>{{end}}
			<form method="POST" action="/signup">

				{{if .SignUpSuccessful}}
					<p>Registration completed. Go to <a href="/login">Login!</a></p>
				{{else}}

					<label for="username">Username:</label>
					<input type="text"  id="username" name="username" value="">
						
					<label for="pass1">Password:</label>
					<input type="password" id="pass1" name="password" value="">
					<input type="submit" name="signup" value="Sign In!">
  				{{end}}
			</form>
		</div>
	</body>
	</html>`

const LoginPage = `<!DOCTYPE html>
	<html>
		<body>
			<div style="display: flex; flex-direction: column; align-items: center" >
			<h1>LOGIN</h1>
	{{if .LoginError}} 		<p style="color:red ; text-transform: uppercase">Either username or password is not in our record! <a href="/signup">Sign Up?</a></p> {{end}}
	{{if .UserNotDatabase}} <p style="color:blue ; text-transform: uppercase">user with such username and password does not exist! <a href="/signup">Sign Up?</a></p> {{end}}

  			<form method="POST" action="/login">
				  {{if .Username}}
					<p><b>{{.Username}}</b>, you're already logged in! <a href="/logout">Logout!</a></p>
				  {{else}}
					<label>Username:</label>
					<input type="text" name="username"><br>
		
					<label>Password:</label>
					<input type="password" name="password">
		
					<input type="submit" name="login" value="Let me in!">
				  {{end}}
  			</form>
			</div>
	  </body>
	</html>`

const DashboardPage = `<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>My profile</title>
	</head>
	<body>
	{{if .UserNotDatabase}} <p style="color:blue ; text-transform: uppercase">user with such username and password does not exist! Sign Up?</p> {{else}}

	<div style="display: flex; flex-direction: column; align-items: center" >
{{if .Logout}}
	
{{else}}
		<p style="color:#2F4F4F; text-transform: uppercase">User:{{.username}}</p>
		<form method="post" action="/logout">
			<button type="submit">Logout</button>
		</form>
{{end}}
	</div>
{{end}}
	</body>
	</html>`

var StartPageTemplate = template.Must(template.New("").Parse(StartPage))
var SignUpTemplate = template.Must(template.New("").Parse(SingUpPage))
var LoginUserTemplate = template.Must(template.New("").Parse(LoginPage))
var DashboardTemplate = template.Must(template.New("").Parse(DashboardPage))
