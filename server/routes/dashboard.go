package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

type dashboardTemplate struct {
	Date string
	User string
	Imgs []imageNames
}

type imageNames struct {
	Name string
}

type UserDataPost struct {
	Uname string    `json:"uname"`
	Pwd   string    `json:"pwd"`
	Data  DataField `json:"data"`
}

type DataField struct {
	X []float64 `json:"x"`
	Y []float64 `json:"y"`
	Z []float64 `json:"z"`
}

func DisplayDashboard(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	sessionToken := cookie.Value
	resp, err := Cache.Do("GET", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if resp == nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	t := time.Now()
	date := fmt.Sprintf("%v %v %v", t.Day(), t.Month(), t.Year())

	data := dashboardTemplate{
		Date: date,
		User: fmt.Sprintf("Welcome %s", resp),
		Imgs: []imageNames{
			{Name: "red.png"},
			{Name: "green.png"},
		},
	}

	tmpl := template.Must(template.ParseFiles(TEMPLATES + "dashboard.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SendUserData(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(uname, pwd)

	pwdCred, ok := credentials[uname]
	if ok && pwd == pwdCred {
		_, err := os.Stat("data/" + uname)
		if os.IsNotExist(err) {
			err = os.Mkdir("data/"+uname, 0755)
		}
		if err == nil {
			var data UserDataPost
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&data)
			if err != nil {
				panic(err)
			}
			fmt.Println(data)
		}
	}
}
