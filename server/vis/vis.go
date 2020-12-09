package vis

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

// chartData is a struct to hold sensor data for line plots
type chartData struct {
	N []int
	A []opts.LineData
	R []opts.LineData
}

// readCSV parses a sensor data csv and returns the parsed data as a chartData
// struct.
func readCSV(filename string) (*chartData, error) {
	// init chart data struct
	var data chartData

	// open the file for parsing
	csvf, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// create a new csv reader
	reader := csv.NewReader(csvf)
	var n int
	var aVals []float64
	var rVals []float64
	for {
		// read the current line of file
		record, err := reader.Read()
		// check for error or at end of file
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// parse accel values in to struct
		a, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			continue
		}

		// parse r value
		r, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}

		// add values to array
		aVals = append(aVals, a)
		rVals = append(rVals, r)
		// add to a slice containing number entries
		data.N = append(data.N, n)
		n++
	}

	// get max value in accel array
	var maxA float64 = math.Inf(-1)
	for _, v := range aVals {
		if v > maxA {
			maxA = v
		}
	}

	// get max value in flex array
	var maxR float64 = math.Inf(-1)
	for _, v := range rVals {
		if v > maxR {
			maxR = v
		}
	}

	// normalize accel and add to data
	for i := range aVals {
		aVals[i] /= maxA
		data.A = append(data.A, opts.LineData{Value: aVals[i], Symbol: "none"})
	}

	// normalize flex and add to data
	for i := range rVals {
		rVals[i] /= maxR
		data.R = append(data.R, opts.LineData{Value: rVals[i], Symbol: "none"})
	}

	return &data, nil
}

// LineChart takes a csv filename and io writer and writes the html and JS
// to the end of the writer
func LineChart(filename string, w io.Writer) error {
	// get data from csv file as chartData struct
	data, err := readCSV(filename)
	if err != nil {
		return err
	}

	// create a new line chart
	line := charts.NewLine()

	fnsplit := strings.Split(filename, "/")[2]
	dateTimeSplit := strings.Split(fnsplit, "T")
	datestring := dateTimeSplit[0]
	timestring := dateTimeSplit[1]
	timestring = strings.Split(timestring, ".")[0]
	timestring = strings.Split(timestring, "-")[0]

	line.SetGlobalOptions(
		// set line chart to desired theme
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemePurplePassion,
		}),
		// set the title of the chart
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("Sensor Data for %s %s", datestring, timestring),
		}),
		// add a legend just above the chart axes
		charts.WithLegendOpts(opts.Legend{
			Show: true,
			Left: "10%",
			Top:  "5%",
		}),
	)

	// add x, y, and z data to chart series
	line.SetXAxis(data.N).
		AddSeries("Acceleration Magnitude", data.A).
		AddSeries("Flex Resistance", data.R)

	// write the chart to the end of the give io writer
	err = line.Render(w)
	if err != nil {
		return err
	}

	return nil
}

// logData is a struct to hold LBP log data
type logData struct {
	Dates []string
	LBP   []float64
	N     []int
}

// readLog parses the lbp data into the log data struct
func readLog(filename string) (*logData, error) {
	var data logData

	csvf, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(csvf)
	var n int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		date := record[0]

		lbp, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}

		data.Dates = append(data.Dates, date)
		data.LBP = append(data.LBP, lbp)
		data.N = append(data.N, n)
		n++
	}

	return &data, nil
}

// addHMData adds a value to the heat map for the given day. -1 is interpreted
// as no data
func addHMData(hmData []opts.HeatMapData, value float64, day int) []opts.HeatMapData {
	if value == -1 {
		hmData = append(hmData, opts.HeatMapData{
			Value: [3]interface{}{
				int((34 - day) % 7), int(day / 7), "-",
			},
		})
	} else {
		hmData = append(hmData, opts.HeatMapData{
			Value: [3]interface{}{
				int((34 - day) % 7), int(day / 7), value,
			},
		})
	}
	return hmData
}

// parseLogData is handles parsing the logData struct into the needed format
// for the heatmap slice.
func parseLogData(data logData) []opts.HeatMapData {
	// make sure that there is data in the passed struct
	if len(data.Dates) < 1 {
		return make([]opts.HeatMapData, 0)
	}
	// get a date from the log to know the month
	date := data.Dates[0]
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		panic(err)
	}

	// create heat map data
	hmData := make([]opts.HeatMapData, 0)

	// fill in days before month start with empty data
	monthStart := t.AddDate(0, 0, -(t.Day() - 1))
	startDay := int(monthStart.Weekday())
	for i := 34; i > (34 - startDay); i-- {
		hmData = addHMData(hmData, -1, i)
	}

	var lbpmean float64
	var entries int
	var lastday int = t.Day()
	var days int = (34 - startDay) - t.Day() + 1
	for i, lbp := range data.LBP {
		// get date of current entry
		date := data.Dates[i]
		t, err := time.Parse(time.RFC3339, date)
		if err != nil {
			panic(err)
		}

		// check if multiple entries for same day //TODO check for skipped days
		if t.Day() == lastday {
			// if there are then add to running average
			lbpmean += lbp
			entries++
		} else {
			// new day so average lbp entries from prev day
			lbpmean /= float64(entries)
			// and append to hm data
			hmData = addHMData(hmData, lbpmean, days)
			// increment day counter
			days--
			// reset number of entries for running avg
			entries = 1
			// and set lbp mean to current entry lbp
			lbpmean = lbp

			// check for skipped days
			if t.Day()-lastday > 1 {
				// and fill in missing data with no value
				for i := lastday; i < t.Day()-1; i++ {
					hmData = addHMData(hmData, -1, i)
					days--
				}
			}
			lastday = t.Day()
		}
	}
	lbpmean /= float64(entries)
	hmData = addHMData(hmData, lbpmean, days)

	return hmData
}

// getCalendarWeeks returns a slice of string with the beginning and start days
// of each week in the given calendar month from date string
func getCalendarWeeks(date string) []string {
	// convert date string to time
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		panic(err)
	}
	// move date to first day of month
	t = t.AddDate(0, 0, -(t.Day() - 1))

	// get weekday of month start
	weekday := t.Weekday()
	// calculate days remaining for first week of month
	daysLeft := 6 - int(weekday)
	// add remaining days
	t = t.AddDate(0, 0, daysLeft)
	var weeks []string
	// append first week
	weeks = append(weeks, fmt.Sprintf("1 - %d", t.Day()))

	// store current month
	month := t.Month()
	// loop over remaining weeks (heatmap always 5x7)
	for i := 0; i < 4; i++ {
		// beginning of next week
		t = t.AddDate(0, 0, 1)
		start := t.Day()
		// 1 week later
		t = t.AddDate(0, 0, 6)
		// make sure didnt go into next month
		if t.Month() != month {
			// if went to next month, minus day to go back to end of prev month
			t = t.AddDate(0, 0, -t.Day())
		}
		// end of week
		end := t.Day()
		weeks = append(weeks, fmt.Sprintf("%d - %d", start, end))
	}
	for i, j := 0, len(weeks)-1; i < j; i, j = i+1, j-1 {
		weeks[i], weeks[j] = weeks[j], weeks[i]
	}

	return weeks
}

// LBPChart writes to the io wrtier the data for the heatmap chart of lbp
func LBPChart(filename string, w io.Writer) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		f.Close()
	}
	// read in the data from the log file
	data, err := readLog(filename)
	if err != nil {
		panic(err)
	}

	// parse the data into the heatmap format
	hmData := parseLogData(*data)

	// create the slice of strings which holds the date ranges for each week
	var weeks []string
	if len(data.Dates) > 0 {
		weeks = getCalendarWeeks(data.Dates[0])
	} else {
		weeks = getCalendarWeeks(time.Now().Format(time.RFC3339))
	}

	// x axis categories are the days of week
	days := [7]string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}

	// get current month for chart title
	month := time.Now().Month().String()

	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		// set color theme
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemePurplePassion,
		}),
		// add title
		charts.WithTitleOpts(opts.Title{
			Title: "Lower Back Pain for " + month,
		}),
		// add x axis labels
		charts.WithXAxisOpts(opts.XAxis{
			Type:      "category",
			Data:      days,
			SplitArea: &opts.SplitArea{Show: true},
		}),
		// add y axis labels
		charts.WithYAxisOpts(opts.YAxis{
			Type:      "category",
			Data:      weeks,
			SplitArea: &opts.SplitArea{Show: true},
		}),
		// enable some chart options
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Min:        0,
			Max:        10,
		}),
	)
	hm.SetXAxis(days).AddSeries("heatmap", hmData)

	err = hm.Render(w)
	if err != nil {
		return err
	}

	return nil
}
