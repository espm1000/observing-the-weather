package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env"
)

type NWSConfig struct {
	BaseURL        string
	GridX          string
	GridY          string
	ForecastOffice string
	StationID      string
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
	LogOutput            string `env:"WEATHER_LOG_FILE" envDefault:"weatherlog.json"`
}

func (n NWSConfig) GetCurrentData() (*CurrentWeatherData, error) {
	var currentData Observation
	slog.Info("getting current weather data", "observationStation", n.StationID, "forecastOffice", n.ForecastOffice)
	resp, err := http.Get(n.BaseURL + "/stations/" + n.StationID + "/observations/latest")
	if err != nil {
		slog.Error("error fetching latest observation data", "error", err)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error reading response body", "error", err)
	}
	slog.Info("response from nws", "response", string(body))
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
	temp_f, err := ConvertCelciusToFahrenheit(currentData.Properties.Temperature.Value)
	if err != nil {
		return nil, err
	}

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
		BaseURL:        "https://api.weather.gov",
		GridX:          "102",
		GridY:          "84",
		ForecastOffice: "MPX",
	}
	slog.Info("getting forecast data")
	resp, err := http.Get(cfg.BaseURL + "/gridpoints/" + cfg.ForecastOffice + "/" + cfg.GridX + "," + cfg.GridY + "/forecast")
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
	logFile, err := os.OpenFile(cfg.LogOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("error creating log file", "error", err)
	}
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	fileLogger := slog.New(slog.NewJSONHandler(logFile, nil))
	slog.SetDefault(fileLogger)
	nws := NWSConfig{
		BaseURL:        "https://api.weather.gov",
		GridX:          "102",
		GridY:          "84",
		ForecastOffice: cfg.ForecastStationId,    // Minneapolis
		StationID:      cfg.ObservationStationId, // St. Paul
	}
	CurrentWeather, err := nws.GetCurrentData()
	if err != nil {
		slog.Error("error", "error", err)
		log.Fatal(err)
	}
	if err := WriteCsv(cfg.ReportOutputDir, *CurrentWeather); err != nil {
		log.Fatal(err)
	}

	PrintToConsole(*CurrentWeather)
}
