package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type dashboardTemplate struct {
	Date string
	User string
	Img  string
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
		Img:  "",
	}

	tmpl := template.Must(template.ParseFiles(TEMPLATES + "dashboard.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
