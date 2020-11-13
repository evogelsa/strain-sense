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
		fmt.Fprintf(w, "User already exists")
	} else if !exist && pwd1 != pwd2 {
		fmt.Fprintf(w, "Passwords don't match'")
	} else if !exist && pwd1 == pwd2 {
		credentials[uname] = pwd1
		fmt.Fprintf(w, "User created")
	}

}

func DisplayCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(TEMPLATES + "create_account.html"))

	data := createTemplate{
		Date: time.Now().Format(time.RFC850)[:17],
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
