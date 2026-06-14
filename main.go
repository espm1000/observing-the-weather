package main

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"github.com/espm1000/observing-the-weather/pkg/observations"
	"github.com/espm1000/observing-the-weather/pkg/report"
	"github.com/espm1000/observing-the-weather/pkg/tools"

	"github.com/caarlos0/env"
)

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

	// PrintToConsole(*CurrentWeather)
}
