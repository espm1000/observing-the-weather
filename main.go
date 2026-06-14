package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/espm1000/observing-the-weather/pkg/observations"
	"github.com/espm1000/observing-the-weather/pkg/report"
	"github.com/espm1000/observing-the-weather/pkg/tools"

	"github.com/caarlos0/env"
)

type NWSConfig struct {
	BaseURL        string
	GridX          string
	GridY          string
	ForecastOffice string
	StationID      string
}

type ForecastWeatherData struct {
	Temperature    any
	Humidity       float64
	Windspeed      any
	ChanceOfPrecip bool
	PrecipPercent  any
	Timestamp      time.Time
}

func (n NWSConfig) GetCurrentData() (*report.CurrentWeatherData, error) {
	var currentData observations.Observation
	slog.Info("getting current weather data", "observationStation", n.StationID, "forecastOffice", n.ForecastOffice)
	resp, err := http.Get(n.BaseURL + "/stations/" + n.StationID + "/observations/latest")
	if err != nil {
		slog.Error("error fetching latest observation data", "error", err)
		return nil, err
	}
	// _, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	slog.Error("error reading response body", "error", err)
	// }

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("error closing stream", "error", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}
	if err = json.NewDecoder(resp.Body).Decode(&currentData); err != nil {
		slog.Error("error decoding response stream", "error", err)
		return nil, err
	}
	temp_f, err := tools.ConvertCelciusToFahrenheit(currentData.Properties.Temperature.Value)
	if err != nil {
		return nil, err
	}

	return &report.CurrentWeatherData{
		Temperature:    temp_f,
		Humidity:       currentData.Properties.RelativeHumidity.Value,
		Windspeed:      currentData.Properties.WindSpeed.Value,
		Timestamp:      currentData.Properties.Timestamp,
		ChanceOfPrecip: false,
	}, nil

}

func PrintToConsole(d report.CurrentWeatherData) {
	fmt.Printf("\nCurrent weather for %v \n", d.Timestamp)
	fmt.Printf("Current Temp: %v F\n", d.Temperature)
	fmt.Printf("Current Windspeed: %v km/h\n", d.Windspeed)
	fmt.Printf("Current Humidity: %v Percent\n", strconv.FormatFloat(d.Humidity, 'f', 2, 64))
	fmt.Printf("Chance of Precip: %v (not implemented)\n", d.ChanceOfPrecip)
}

func main() {
	cfg := tools.Environment{}
	if err := env.Parse(&cfg); err != nil {
		slog.Error("failed to parse env vars")
	}
	logger, err := tools.SetLogger(cfg)
	if err != nil {
		slog.Error("error setting logger", "error", err)
		panic(err)
	}
	slog.SetDefault(logger)
	nws := NWSConfig{
		BaseURL:        "https://api.weather.gov",
		GridX:          "102",
		GridY:          "84",
		ForecastOffice: cfg.ForecastStationId,    // Minneapolis
		StationID:      cfg.ObservationStationId, // St. Paul
	}
	CurrentWeather, err := nws.GetCurrentData()
	if err != nil {
		slog.Error("error getting weather", "error", err)
		log.Fatal(err)
	}
	if err := report.WriteCsv(cfg.ReportOutputDir, *CurrentWeather); err != nil {
		slog.Error("error writing csv", "error", err)
		log.Fatal(err)
	}

	// PrintToConsole(*CurrentWeather)
}
