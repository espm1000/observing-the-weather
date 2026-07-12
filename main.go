package main

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	"github.com/espm1000/observing-the-weather/pkg/client"
	"github.com/espm1000/observing-the-weather/pkg/nws"
	"github.com/espm1000/observing-the-weather/pkg/report"
	"github.com/espm1000/observing-the-weather/pkg/tools"
)

func setPreConfig() (*tools.Environment, *report.ReportConfig, error) {
	rpt := report.ReportConfig{}
	cfg := tools.Environment{}
	if err := tools.SetEnvironment(&cfg); err != nil {
		slog.Error("error setting environment variables", "error", err)
		return nil, nil, err
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
	slog.Debug("report config", "directory", rpt.Directory, "reportFile", rpt.NWSReport)

	return &cfg, &rpt, err
}

func setReportEnvironment(r *report.ReportConfig) error {
	if err := env.Parse(r); err != nil {
		return err
	}
	return nil
}

func main() {
	env, rptCfg, _ := setPreConfig()
	if err := NWS(env, rptCfg); err != nil {
		panic(err)
	}
	if err := getNCEIWeather(env.NCEIToken, rptCfg); err != nil {
		slog.Error("error getting ncei data", "error", err)
	}
	startWebServer()
}

func getNCEIWeather(token string, r *report.ReportConfig) error {
	err := report.NCEIReport(report.Params{
		StartDate:       "2025-07-08",
		EndDate:         "2026-07-08",
		Units:           "standard",
		Dataset:         "GHCND",
		StationId:       "USW00014922",
		Limit:           "1000",
		IncludeMetadata: "false",
	}, token, r)
	if err != nil {
		return err
	}

	return nil
}

func NWS(e *tools.Environment, rc *report.ReportConfig) error {
	if e.NWSEnabled == "false" {
		slog.Info("nws disabled by environment variable")
		return nil
	}
	httpCfg := client.HttpClientConfig{
		UserAgent: "weather-app@esp.m1k@gmail.com",
		Timeout:   10 * time.Second,
	}
	nws := nws.NWSConfig{
		StationID: e.ObservationStationId, // St. Paul
	}
	CurrentWeather, err := nws.GetCurrentData(&httpCfg)
	if err != nil {
		slog.Error("error getting weather", "error", err)
		return err
	}
	if err := report.WriteCsv(*rc, *CurrentWeather); err != nil {
		slog.Error("error writing csv", "error", err)
		return err
	}
	if e.PrintToConsole == "true" {
		tools.PrintToConsole(*CurrentWeather, *e)
	}
	return err
}

func startWebServer() error {
	log.Println("starting http server...")
	http.HandleFunc("/data/", report.CsvHandler("ncei.csv"))
	if err := http.ListenAndServe(":5050", nil); err != nil {
		log.Fatal("server failed to start: ", err)
		return err
	}
	return nil
}
