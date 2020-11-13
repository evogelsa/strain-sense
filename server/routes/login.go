package routes

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// credentials is map from [user] to pwd
var credentials map[string]string
var loadedCreds bool

type loginTemplate struct {
	Date string
}

func initCredentials() {
	_, err := os.Stat("creds")
	if err == nil {
		// exists
		loadCredentials()
	} else {
		credentials = make(map[string]string)
	}
}

func loadCredentials() {
	f, err := os.Open("creds")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	err = dec.Decode(&credentials)
	if err != nil {
		log.Fatal(err)
	}
}

func saveCredentials() {
	f, err := os.Create("creds")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(credentials)
	if err != nil {
		log.Fatal(err)
	}
}

func AuthenticateLogin(w http.ResponseWriter, r *http.Request) {
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
	pwd := r.PostForm.Get("pwd")

	pwdCred, ok := credentials[uname]
	if ok && pwd == pwdCred {
		http.Redirect(w, r, "/wearables/dashboard", http.StatusSeeOther)
	} else {
		DisplayLogin(w, r)
	}
}

func DisplayLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(TEMPLATES + "login.html"))

	data := loginTemplate{
		Date: time.Now().Format(time.RFC850)[:17],
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
