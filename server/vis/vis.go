package vis

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type chartData struct {
	N []int
	X []opts.LineData
	Y []opts.LineData
	Z []opts.LineData
}

func readCSV(filename string) (*chartData, error) {
	var data chartData

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

		x, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, err
		}

		y, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}

		z, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}

		data.X = append(data.X, opts.LineData{Value: x})
		data.Y = append(data.Y, opts.LineData{Value: y})
		data.Z = append(data.Z, opts.LineData{Value: z})
		data.N = append(data.N, n)
		n++
	}

	return &data, nil
}

func LineChart(filename string, w io.Writer) error {
	data, err := readCSV(filename)
	if err != nil {
		return err
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(
			opts.Initialization{Theme: types.ThemePurplePassion},
		),
		charts.WithTitleOpts(opts.Title{
			Title: filename,
		}),
	)

	line.SetXAxis(data.N).
		AddSeries("X Acceleration", data.X).
		AddSeries("X Acceleration", data.Y).
		AddSeries("X Acceleration", data.Z)

	err = line.Render(w)
	if err != nil {
		return err
	}

	return nil
}
