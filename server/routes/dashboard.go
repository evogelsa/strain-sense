package routes

import (
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
	data := dashboardTemplate{
		Date: time.Now().Format(time.RFC850)[:17],
		User: r.PostForm.Get("uname"),
		Img:  "",
	}
	tmpl := template.Must(template.ParseFiles(TEMPLATES + "dashboard.html"))

	err := tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
