package main

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"github.com/espm1000/observing-the-weather/pkg/observations"
	"github.com/espm1000/observing-the-weather/pkg/report"
	"github.com/espm1000/observing-the-weather/pkg/tools"
)

func PrintToConsole(d report.CurrentWeatherData) {
	fmt.Printf("\nCurrent weather for %v \n", d.Timestamp)
	fmt.Printf("Current Temp: %v F\n", strconv.FormatFloat(d.Temperature, 'f', 2, 64))
	fmt.Printf("Current Windspeed: %v km/h\n", d.Windspeed)
	fmt.Printf("Current Humidity: %v Percent\n", strconv.FormatFloat(d.Humidity, 'f', 2, 64))
	fmt.Printf("Chance of Precip: %v (not implemented)\n", d.ChanceOfPrecip)
}

func setPreConfig() *tools.Environment {
	cfg := tools.Environment{}
	if err := tools.SetEnvironment(&cfg); err != nil {
		slog.Error("error setting environment variables", "error", err)
		return nil
	}
	logger, err := tools.SetLogger(cfg)
	if err != nil {
		slog.Error("error setting logger", "error", err)
	}
	cfg.Logger = logger

	return &cfg
}

func main() {
	cfg := setPreConfig()
	slog.SetDefault(cfg.Logger)
	nws := observations.NWSConfig{
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
	//PrintToConsole(*CurrentWeather)
}
