package ncei

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	nceiBaseUrl      = "https://www.ncei.noaa.gov/cdo-web/api/v2"
	dataEndpoint     = "data"
	dataSetEnpoint   = "datasets"
	dataTypeEndpoint = "datatypes"
	locationEndpoint = "locations"
	stationEndpoint  = "stations"
)

type Params struct {
	StartDate string
	EndDate   string
	StationId string
	Dataset   string
	Units     string
	Format    string
	Limit     string
}

type WeatherData struct {
	Date     string `json:"date"`
	Location string `json:"location"`
	MaxTemp  string `json:"maxTemp"`
	MinTemp  string `json:"minTemp"`
	Precip   bool   `json:"precip"`
}

type Report struct {
	header   []string
	datarows WeatherData
}

func buildReportFile(w *WeatherData) *Report {
	headers := []string{"date", "location", "maxTemp", "minTemp", "precip"}

	return &Report{
		header: headers,
		datarows: WeatherData{
			Date:     w.Date,
			Location: w.Location,
			MaxTemp:  w.MaxTemp,
			MinTemp:  w.MinTemp,
			Precip:   w.Precip,
		},
	}
}

func getRawData(p Params, nceiToken string) (string, error) {
	params := url.Values{}
	params.Add("startdate", p.StartDate)
	params.Add("endDate", p.EndDate)
	params.Add("datasetid", p.Dataset)
	params.Add("stationId", p.Dataset+":"+p.StationId)
	params.Add("format", p.Format)
	params.Add("units", p.Units)
	params.Add("limit", p.Limit)

	url := nceiBaseUrl + "?" + params.Encode()

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
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v\n", err)
	}
	log.Println(string(respBytes))
	return "", nil
}

func NewReport(p Params, token string) (*WeatherData, error) {
	_, err := getRawData(p, token)
	if err != nil {
		return nil, err
	}

	return &WeatherData{
		Date:     p.StartDate + "===>" + p.EndDate,
		Location: p.StationId,
		MaxTemp:  "hot",
		MinTemp:  "cold",
		Precip:   true,
	}, nil
}
