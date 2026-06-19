package main

import (
	"log/slog"

	"github.com/caarlos0/env"
	"github.com/espm1000/observing-the-weather/pkg/nws"
	"github.com/espm1000/observing-the-weather/pkg/report"
	"github.com/espm1000/observing-the-weather/pkg/tools"
)

func setPreConfig() (*tools.Environment, *report.ReportConfig) {
	rpt := report.ReportConfig{}
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
	if err := setReportEnvironment(&rpt); err != nil {
		slog.Error("error setting report folders", "error", err)
	}
	slog.Debug("report config", "directory", rpt.Directory, "reportFile", rpt.ReportFile)

	return &cfg, &rpt
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
		tools.PrintToConsole(*CurrentWeather, *cfg)
	}
	return err
}

func setReportEnvironment(r *report.ReportConfig) error {
	if err := env.Parse(r); err != nil {
		return err
	}
	return nil
}
