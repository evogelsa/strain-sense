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

// Cache is the redis server connection
var Cache redis.Conn

// loginTemplate is a struct to hold data for the login.html file template
type loginTemplate struct {
	Date string
}

// InitInitCredentials checks if a creds file exists and if so loads it into the
// credentials memory. Otherwise it will generate a new credentials map
func InitCredentials() {
	_, err := os.Stat("creds")
	if err == nil {
		// exists
		loadCredentials()
	} else {
		credentials = make(map[string]string)
	}
}

// loadCredentials should only be called by InitCredentials. It handles opening
// the creds file and decoding it into the credentials map
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

// saveCredentials saves the credentials map as a gob (not cryto secure)
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

// AAuthenticateLogin handles the login POST routine and validates that the
// submitted uname as pwd are valid and match a record in the credentials map.
func AuthenticateLogin(w http.ResponseWriter, r *http.Request) {
	// parse login submission form
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// retrieve uname and pwd from html form
	uname := r.PostForm.Get("uname")
	pwd := r.PostForm.Get("pwd")

	// get existing password from user and check that user exists
	pwdCred, ok := credentials[uname]
	if ok && pwd == pwdCred {
		// if user is existant and the password was correct then
		// create session token and add to cache
		sessionToken := uuid.NewV4().String()
		// token expires after 60 minutes and stores username
		_, err = Cache.Do("SETEX", sessionToken, "3600", uname)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// set the http cookie to the freshly generated session token cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(time.Hour),
		})

		// redirect user upon success to the dashboard
		http.Redirect(w, r, "/wearables/dashboard", http.StatusSeeOther)
	} else {
		// if user could not be validated refresh login page
		DisplayLogin(w, r)
	}
}

// DisplayLogin handles the GET routine for the login page. It shows a basic
// login form
func DisplayLogin(w http.ResponseWriter, r *http.Request) {
	// load login template
	tmpl := template.Must(template.ParseFiles(TEMPLATES + "login.html"))

	// fetch current data and add to login template
	t := time.Now()
	date := fmt.Sprintf("%v %v %v", t.Day(), t.Month(), t.Year())
	data := loginTemplate{
		Date: date,
	}

	// display the login page
	err := tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
