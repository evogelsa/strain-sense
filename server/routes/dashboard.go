package routes

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.iu.edu/evogelsa/strain-sense/vis"
)

type dashboardTemplate struct {
	Date string
	User string
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
	}

	tmpl := template.Must(template.ParseFiles(TEMPLATES + "dashboard.html"))

	var file bytes.Buffer
	err = tmpl.Execute(&file, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	dirname := fmt.Sprintf("data/%s/", resp)
	dir, err := os.Open(dirname)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	filenames, err := dir.Readdirnames(0)
	if err != nil {
		panic(err)
	}

	var bodies []string

	for _, filename := range filenames {
		var buf bytes.Buffer
		err = vis.LineChart(dirname+string(filename), &buf)
		if err != nil {
			panic(err)
		}

		doc, err := goquery.NewDocumentFromReader(&buf)
		if err != nil {
			panic(err)
		}

		doc.Find("body").Each(func(i int, s *goquery.Selection) {
			body, _ := s.Html()
			bodies = append(bodies, `<div>`+body+`</div><br>`)
		})
	}

	for _, body := range bodies {
		file.WriteString(body)
	}

	fmt.Fprint(w, file.String())
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

		var filename string
		if err == nil {
			filename = "data/" + data.Uname + "/" +
				time.Now().Format(time.RFC3339) + ".csv"
			csvf, err := os.Create(filename)
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
