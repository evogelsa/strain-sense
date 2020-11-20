package routes

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
)

// credentials is map from [user] to pwd
var credentials map[string]string
var Cache redis.Conn

type loginTemplate struct {
	Date string
}

func InitCredentials() {
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
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	uname := r.PostForm.Get("uname")
	pwd := r.PostForm.Get("pwd")

	pwdCred, ok := credentials[uname]
	if ok && pwd == pwdCred {
		// create session token and add to cache
		sessionToken := uuid.NewV4().String()
		_, err = Cache.Do("SETEX", sessionToken, "3600", uname)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(time.Hour),
		})

		http.Redirect(w, r, "/wearables/dashboard", http.StatusSeeOther)
	} else {
		DisplayLogin(w, r)
	}
}

func DisplayLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(TEMPLATES + "login.html"))

	t := time.Now()
	date := fmt.Sprintf("%v %v %v", t.Day(), t.Month(), t.Year())
	data := loginTemplate{
		Date: date,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
