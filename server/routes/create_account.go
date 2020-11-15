package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type createTemplate struct {
	Date string
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	if !loadedCreds {
		initCredentials()
		loadedCreds = true
	} else {
		saveCredentials()
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	uname := r.PostForm.Get("uname")
	pwd1 := r.PostForm.Get("psw1")
	pwd2 := r.PostForm.Get("psw2")

	_, exist := credentials[uname]
	if exist {
		tmpl := template.Must(template.ParseFiles(
			TEMPLATES + "user_exists.html",
		))
		err := tmpl.Execute(w, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if !exist && pwd1 != pwd2 {
		tmpl := template.Must(template.ParseFiles(
			TEMPLATES + "password_no_match.html",
		))
		err := tmpl.Execute(w, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if !exist && pwd1 == pwd2 {
		credentials[uname] = pwd1
		tmpl := template.Must(template.ParseFiles(
			TEMPLATES + "user_created.html",
		))
		err := tmpl.Execute(w, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func DisplayCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		TEMPLATES + "create_account.html",
	))

	t := time.Now()
	date := fmt.Sprintf("%v %v %v", t.Day(), t.Month(), t.Year())
	data := createTemplate{
		Date: date,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
