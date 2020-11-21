package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

// createTemplate is a struct to hold data for the create_account.html template
// file
type createTemplate struct {
	Date string
}

// CreateAccount is the POST routine for create_account and handles validating
// the submitted user data and adding to credentials. On sucess will display
// a success message on the server.
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	// parse the HTML form data
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// retrieve username and passwords from form
	uname := r.PostForm.Get("uname")
	pwd1 := r.PostForm.Get("pwd1")
	pwd2 := r.PostForm.Get("pwd2")

	// check if the given username is already existant
	_, exist := credentials[uname]
	if exist {
		// if user is already taken display error
		tmpl := template.Must(template.ParseFiles(
			TEMPLATES + "user_exists.html",
		))
		err := tmpl.Execute(w, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if !exist && pwd1 != pwd2 {
		// if user is not taken but passwords don't match display error
		tmpl := template.Must(template.ParseFiles(
			TEMPLATES + "password_no_match.html",
		))
		err := tmpl.Execute(w, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if !exist && pwd1 == pwd2 {
		// if user doesn't exist and passwords match create user account
		credentials[uname] = pwd1
		saveCredentials()

		// create a data directory to store user data
		err := os.MkdirAll("data/"+uname+"/", 0755)
		if err != nil {
			panic(err)
		}

		// display success message
		tmpl := template.Must(template.ParseFiles(
			TEMPLATES + "user_created.html",
		))
		err = tmpl.Execute(w, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// DisplayCreate displays the user creation page
func DisplayCreate(w http.ResponseWriter, r *http.Request) {
	// parse creation page html
	tmpl := template.Must(template.ParseFiles(
		TEMPLATES + "create_account.html",
	))

	// get tume and add to template data
	t := time.Now()
	date := fmt.Sprintf("%v %v %v", t.Day(), t.Month(), t.Year())
	data := createTemplate{
		Date: date,
	}

	// show page and update date with server date
	err := tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
