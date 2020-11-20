package routes

import (
	"encoding/csv"
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
	Uname string      `json:"uname"`
	Pwd   string      `json:"pwd"`
	Data  []DataField `json:"data"`
}

type DataField struct {
	X string `json:"X"`
	Y string `json:"Y"`
	Z string `json:"Z"`
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
	var data UserDataPost
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	pwdCred, ok := credentials[data.Uname]
	if ok && data.Pwd == pwdCred {
		_, err := os.Stat("data/" + data.Uname)
		if os.IsNotExist(err) {
			err = os.MkdirAll("data/"+data.Uname, 0755)
		}

		if err == nil {
			csvf, err := os.Create("data/" + data.Uname + "/" +
				time.Now().Format(time.RFC3339) + ".csv")
			if err != nil {
				panic(err)
			}
			defer csvf.Close()

			writer := csv.NewWriter(csvf)
			for _, xyz := range data.Data {
				var row []string
				row = append(row, xyz.X)
				row = append(row, xyz.Y)
				row = append(row, xyz.Z)
				err = writer.Write(row)
				if err != nil {
					panic(err)
				}
			}
			writer.Flush()
		} else {
			panic(err)
		}
	}
}
