package report

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/espm1000/observing-the-weather/pkg/ncei"
)

type Params struct {
	StartDate       string
	EndDate         string
	StationId       string
	Dataset         string
	Units           string
	Limit           string
	IncludeMetadata string
}

type Report struct {
	header   []string
	datarows ncei.WeatherData
}

func NCEIReport(p Params, token string) (*Report, error) {
	log.Println("generating ncei weather report")
	if token == "" {
		return nil, errors.New("no ncei token")
	}
	rd, err := getRawData(p, token)
	if err != nil {
		return nil, err
	}

	report, err := buildReport(rd)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func BuildNCEI(token string) *ncei.Config {
	return ncei.New(token)
}

func getRawData(p Params, nceiToken string) (*ncei.WeatherResponse, error) {
	cfg := BuildNCEI(nceiToken)

	params := url.Values{}
	params.Add("startdate", p.StartDate)
	params.Add("enddate", p.EndDate)
	params.Add("datasetid", p.Dataset)
	params.Add("stationid", p.Dataset+":"+p.StationId)
	params.Add("units", p.Units)
	params.Add("limit", p.Limit)
	params.Add("includemetadata", p.IncludeMetadata)

	url := cfg.BaseURL + cfg.DataEP + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("error generating http request: ", err)
	}

	req.Header.Add("token", nceiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error fetching ncei data: ", err)
	}

	defer resp.Body.Close()

	rawWeather := ncei.WeatherResponse{}

	err = json.NewDecoder(resp.Body).Decode(&rawWeather)
	if err != nil {
		log.Println("error decoding json: ", err)
		return nil, err
	}

	return &rawWeather, nil
}

func buildReport(w *ncei.WeatherResponse) (*Report, error) {
	header := []string{"date", "location", "maxTemp", "minTemp", "precip"}
	report := Report{
		header:   header,
		datarows: ncei.WeatherData{},
	}

	file, err := os.Create("ncei.csv")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		log.Println("error writing header: ", err)
		return nil, err
	}

	// AI Assisted: Getting all data from the response and converting it to the WeatherData struct in a loop.
	// Learned to make a map with using the struct within the map.
	weatherByDate := make(map[string]*ncei.WeatherData)
	for _, data := range w.Results {
		wd, ok := weatherByDate[data.Date]
		if !ok {
			wd = &ncei.WeatherData{Date: data.Date, Location: data.Station}
			weatherByDate[data.Date] = wd
		}

		switch data.Datatype {
		case "TMAX":
			wd.MaxTemp = convertFtoS(data.Value)
		case "TMIN":
			wd.MinTemp = convertFtoS(data.Value)
		case "PRCP":
			wd.Precip = convertFtoS(data.Value) != "0" && convertFtoS(data.Value) != ""
		}
	}

	for _, i := range weatherByDate {
		row := []string{
			i.Date,
			i.Location,
			i.MaxTemp,
			i.MinTemp,
			"null",
		}
		if err := writer.Write(row); err != nil {
			log.Println("error writing row data: ", err)
		}
	}

	return &report, nil
}

func convertFtoS(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64)
}
