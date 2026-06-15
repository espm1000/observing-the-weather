package main

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/espm1000/observing-the-weather/pkg/nws"
	"github.com/espm1000/observing-the-weather/pkg/report"
	"github.com/espm1000/observing-the-weather/pkg/tools"
)

func PrintToConsole(d report.CurrentWeatherData, cfg tools.Environment) {
	fmt.Printf("\nCurrent weather for %v \n", d.Timestamp)
	fmt.Printf("Current Temp: %v F\n", strconv.FormatFloat(d.Temperature, 'f', 2, 64))
	fmt.Printf("Current Windspeed: %v km/h\n", d.Windspeed)
	fmt.Printf("Current Humidity: %v Percent\n", strconv.FormatFloat(d.Humidity, 'f', 2, 64))
	fmt.Printf("Grid Coordinates\nX: %v\nY: %v\n", cfg.GridX, cfg.GridY)
	fmt.Printf("Chance of Precip: %v (not implemented)\n", d.ChanceOfPrecip)
}

func setPreConfig() (*tools.Environment, *report.ReportConfig) {
	cfg := tools.Environment{}
	if err := tools.SetEnvironment(&cfg); err != nil {
		slog.Error("error setting environment variables", "error", err)
		return nil, nil
	}
	logger, err := tools.SetLogger(cfg)
	if err != nil {
		slog.Error("error setting logger", "error", err)
	}
	cfg.Logger = logger
	slog.SetDefault(cfg.Logger)

	rpt := setReportConfig()
	slog.Debug("report config", "directory", rpt.Directory, "reportFile", rpt.ReportFile)

	return &cfg, &rpt
}

func setReportConfig() report.ReportConfig {
	rpt := report.ReportConfig{}
	if err := tools.SetReportEnvironment(&rpt); err != nil {
		slog.Error("error setting report vars", "error", err)
	}
	return rpt
}

func main() {
	if err := Main(); err != nil {
		panic(err)
	}
}

func Main() error {
	cfg, rpt := setPreConfig()
	nws := nws.NWSConfig{
		StationID: cfg.ObservationStationId, // St. Paul
	}
	CurrentWeather, err := nws.GetCurrentData()
	if err != nil {
		slog.Error("error getting weather", "error", err)
		return err
	}
	if err := report.WriteCsv(*rpt, *CurrentWeather); err != nil {
		slog.Error("error writing csv", "error", err)
		return err
	}
	if cfg.PrintToConsole == "true" {
		PrintToConsole(*CurrentWeather, *cfg)
	}
	return err
}
