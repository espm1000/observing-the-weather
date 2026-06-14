package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env"
)

const (
	BaseURL       = "https://api.weather.gov"
	HistoricalURI = "/stations/%s/observations"
)

type NWSConfig struct {
	GridX          string
	GridY          string
	ForecastOffice string
	StationID      string
	StartDate      string
	EndDate        string
}

type CurrentWeatherData struct {
	Temperature    float64
	Humidity       float64
	Windspeed      any
	ChanceOfPrecip bool
	Timestamp      string
	PrecipLastHour float64
}

type ForecastWeatherData struct {
	Temperature    any
	Humidity       float64
	Windspeed      any
	ChanceOfPrecip bool
	PrecipPercent  any
	Timestamp      time.Time
}

type Environment struct {
	ReportOutputDir      string `env:"WEATHER_REPORT_DIR" envDefault:"/data"`
	ObservationStationId string `env:"WEATHER_OBSERVATION_STATION_ID" envDefault:"KSTP"`
	ForecastStationId    string `env:"WEATHER_FORECAST_STATION_ID" envDefault:"MPX"`
	LogDirectory         string `env:"WEATHER_LOG_DIRECTORY" envDefault:"logs"`
	LogOutput            string `env:"WEATHER_LOG_FILE" envDefault:"weatherlog.json"`
	LogLevel             slog.Level
}

func (n NWSConfig) GetCurrentData() (*CurrentWeatherData, error) {
	var currentData Observation
	slog.Info("getting current weather data", "observationStation", n.StationID, "forecastOffice", n.ForecastOffice)
	resp, err := http.Get(BaseURL + "/stations/" + n.StationID + "/observations/latest")
	if err != nil {
		slog.Error("error fetching latest observation data", "error", err)
		return nil, err
	}

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
	temp_f := ConvertCelciusToFahrenheit(currentData.Properties.Temperature.Value)
	
	return &CurrentWeatherData{
		Temperature:    temp_f,
		Humidity:       currentData.Properties.RelativeHumidity.Value,
		Windspeed:      currentData.Properties.WindSpeed.Value,
		Timestamp:      currentData.Properties.Timestamp,
		ChanceOfPrecip: false,
	}, nil

}

func (n NWSConfig) GetForecastData() (*ForecastWeatherData, error) {
	var result Forecast

	cfg := NWSConfig{

		GridX:          "102",
		GridY:          "84",
		ForecastOffice: "MPX",
	}
	slog.Info("getting forecast data")
	resp, err := http.Get(BaseURL + "/gridpoints/" + cfg.ForecastOffice + "/" + cfg.GridX + "," + cfg.GridY + "/forecast")
	if err != nil {
		slog.Error("error calling weather service api", "error", err)
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("error closing stream", "error", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		slog.Error("non-200 response from upstread", "status_code", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		slog.Error("failed to read response body", "error", err)
		return nil, err
	}
	return &ForecastWeatherData{
		Temperature:   result.Properties.Periods[0].Temperature,
		Windspeed:     result.Properties.Periods[0].WindSpeed,
		Timestamp:     result.Properties.Periods[0].EndTime,
		PrecipPercent: result.Properties.Periods[0].ProbabilityOfPrecipitation.Value,
	}, nil
}

func PrintToConsole(d CurrentWeatherData) {
	fmt.Printf("\nCurrent weather for %v \n", d.Timestamp)
	fmt.Printf("Current Temp: %v F\n", d.Temperature)
	fmt.Printf("Current Windspeed: %v km/h\n", d.Windspeed)
	fmt.Printf("Current Humidity: %v Percent\n", strconv.FormatFloat(d.Humidity, 'f', 2, 64))
	fmt.Printf("Chance of Precip: %v (not implemented)\n", d.ChanceOfPrecip)
}

func main() {
	cfg := Environment{}
	if err := env.Parse(&cfg); err != nil {
		slog.Error("failed to parse env vars")
	}
	logger, err := SetLogger(cfg)
	if err != nil {
		slog.Error("error setting logger", "error", err)
		panic(err)
	}
	slog.SetDefault(logger)
	nws := NWSConfig{
		GridX:          "102",
		GridY:          "84",
		ForecastOffice: cfg.ForecastStationId,    // Minneapolis
		StationID:      cfg.ObservationStationId, // St. Paul
		StartDate:      "2026-06-01T00:00:00Z",
		EndDate:        "2026-06-10T00:00:00Z",
	}
	CurrentWeather, err := nws.GetCurrentData()
	if err != nil {
		slog.Error("error getting weather", "error", err)
		log.Fatal(err)
	}
	if err := WriteCsv(cfg.ReportOutputDir, *CurrentWeather); err != nil {
		slog.Error("error writing csv", "error", err)
		log.Fatal(err)
	}

	PrintToConsole(*CurrentWeather)
}
