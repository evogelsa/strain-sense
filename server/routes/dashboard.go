package routes

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.iu.edu/evogelsa/strain-sense/vis"
)

// dashboardTemplate is a struct to hold the data for replacing the dashboard
// html file
type dashboardTemplate struct {
	Date string
	User string
}

// UserDataPost is a struct that is used to unmarshal a json POST request into.
// It stores the POST request body and can be saved as a csv.
type UserDataPost struct {
	Uname string      `json:"uname"`
	Pwd   string      `json:"pwd"`
	Data  []DataField `json:"data"`
}

// DataField is used inside UserDataPost and stores the actual sensor data from
// the post request
type DataField struct {
	A string `json:"A"`
	R string `json:"R"`
}

// DisplayDashboard shows the dashboard webpage and validates the user
// attempting to access the webpage. Will redirect to login if no cookie or
// cookie invalid.
func DisplayDashboard(w http.ResponseWriter, r *http.Request) {
	// retrieve session cookie and check for existance
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		// if no cookie user not authorized, redirect to login
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/wearables/login", http.StatusSeeOther)
	} else if err != nil {
		// if other error bad request and redirect to login
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/wearables/login", http.StatusSeeOther)
	}

	// if cookie valid get session token from the cookie
	sessionToken := cookie.Value
	// resp will be the username from the cookie
	resp, err := Cache.Do("GET", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if resp == nil {
		// if resp is nil then user not authorized, redirect to login
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/wearables/login", http.StatusSeeOther)
	}

	// get current date to display at top of webpage
	t := time.Now()
	date := fmt.Sprintf("%v %v %v", t.Day(), t.Month(), t.Year())

	// update dashboard template data with date and user name
	data := dashboardTemplate{
		Date: date,
		User: fmt.Sprintf("Welcome %s", resp),
	}

	// parse the dashboard html template
	tmpl := template.Must(template.ParseFiles(TEMPLATES + "dashboard.html"))

	// create a bytes buffer to write template to. use of a buffer allows adding
	// extra writes after the template (used for showing the charts)
	var file bytes.Buffer
	// write template data to buffer
	err = tmpl.Execute(&file, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// open the data directory of the user
	dirname := fmt.Sprintf("data/%s/", resp)
	dir, err := os.Open(dirname)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	// each file in the data directory contains some set of data so load those
	// filenames to access later
	filenames, err := dir.Readdirnames(0)
	if err != nil {
		panic(err)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(filenames)))

	// create a slice which will hold the html and js for the charts
	var bodies []string

	// vis lbp log first
	var buf bytes.Buffer
	logname := fmt.Sprintf(
		"%s%d-%d_%s",
		dirname,
		t.Year(),
		t.Month(),
		"LBP_log.csv",
	)
	err = vis.LBPChart(logname, &buf)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		panic(err)
	}
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		body, _ := s.Html()
		bodies = append(bodies, body)
	})

	// for each file name in the data directory
	for _, filename := range filenames {
		// skip lbp logs
		if strings.Contains(filename, "LBP_log") {
			continue
		}
		// create a new byte buffer
		var buf bytes.Buffer
		// and read the data from the file and generate a line chart with
		// echarts. LineChart returns the html and js necessesary to display the
		// data as a line chart.
		err = vis.LineChart(dirname+string(filename), &buf)
		if err != nil {
			panic(err)
		}

		// however this data assumes that its the sole element on the page, so
		// parse the html for just the html body and js script
		doc, err := goquery.NewDocumentFromReader(&buf)
		if err != nil {
			panic(err)
		}

		// find the body tag and append the content of the body to the bodies
		// slice made earlier
		doc.Find("body").Each(func(i int, s *goquery.Selection) {
			body, _ := s.Html()
			bodies = append(bodies, body)
		})
	}

	// now append each entry in bodies to the end of the file bytes buffer, this
	// adds the chart to the end of the page
	for _, body := range bodies {
		file.WriteString(body)
	}

	// finally display the generated webpage
	fmt.Fprint(w, file.String())
}

// SendUserData handles the POST routine of sending user data from the device.
// csv data from the device is encoded as json and sent as a POST routine which
// is then decoded back to csv and saved as a file
func SendUserData(w http.ResponseWriter, r *http.Request) {
	// create empty variable struct to decode data to
	var data UserDataPost
	// decode the JSON data into the go struct UserDataPost
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	// retrieve user password from credentials and validate that user exists
	// and password is correct
	pwdCred, ok := credentials[data.Uname]
	if ok && data.Pwd == pwdCred {
		// data is saved as file in the user data directory with the filename
		// set to the datetime it was received in ISO8601 standard
		// set filename to user data directory plus iso8601 time
		filename := "data/" + data.Uname + "/" +
			time.Now().Format(time.RFC3339) + ".csv"

		// create the file
		csvf, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer csvf.Close()

		// make a new writer and decode the UserDataPost struct into a csv
		// format and write to file
		writer := csv.NewWriter(csvf)
		for _, ar := range data.Data {
			var row []string
			row = append(row, ar.A)
			row = append(row, ar.R)
			err = writer.Write(row)
			if err != nil {
				panic(err)
			}
		}
		// flush the writer when done
		writer.Flush()
	}
}

func LBPLog(w http.ResponseWriter, r *http.Request) {
	// retrieve session cookie and check for existance
	cookie, err := r.Cookie("session_token")
	if err == http.ErrNoCookie {
		// if no cookie user not authorized, redirect to login
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/wearables/login", http.StatusSeeOther)
	} else if err != nil {
		// if other error bad request and redirect to login
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/wearables/login", http.StatusSeeOther)
	}

	// if cookie valid get session token from the cookie
	sessionToken := cookie.Value
	// resp will be the username from the cookie
	resp, err := Cache.Do("GET", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if resp == nil {
		// if resp is nil then user not authorized, redirect to login
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/wearables/login", http.StatusSeeOther)
	}

	uname := fmt.Sprintf("%s", resp)

	t := time.Now()
	dirname := "data/" + uname + "/"
	logname := fmt.Sprintf(
		"%s%d-%d_%s",
		dirname,
		t.Year(),
		t.Month(),
		"LBP_log.csv",
	)
	f, err := os.OpenFile(
		logname,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0644,
	)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	rating := r.PostForm.Get("rating")
	date := time.Now().Format(time.RFC3339)

	var line string = date + "," + rating + "\n"
	_, err = f.WriteString(line)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/wearables/dashboard", http.StatusSeeOther)
}
