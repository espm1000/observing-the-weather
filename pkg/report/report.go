package report

import (
	"encoding/csv"
	"errors"
	"log/slog"
	"os"
	"path"
	"strconv"
	"time"
)

type CurrentWeatherData struct {
	Temperature    float64
	Humidity       float64
	Windspeed      any
	ChanceOfPrecip bool
	Timestamp      string
	PrecipLastHour float64
}

func InitCsv(dir string) error {
	slog.Debug("initializing empty csv report")
	headers := []string{"timestamp", "temperature", "humidity", "precipchance", "polledTimestamp"}
	_, err := os.Stat(path.Join(dir, "currentWeather.csv"))
	if err == nil {
		slog.Debug("report file exists")
		return err
	}
	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err = os.Mkdir(dir, 0755); err != nil {
				slog.Error("error creating directory", "directory", dir, "error", err)
				return err
			}
		}
	}
	file, err := os.Create(path.Join(dir, "currentWeather.csv"))
	if err != nil {
		slog.Error("error creating current report file", "error", err)
		return err
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write(headers); err != nil {
		slog.Error("error writing headers to report", "error", err)
		return err
	}
	return err
}

func WriteCsv(dir string, d CurrentWeatherData) error {
	var reportData []CurrentWeatherData
	if err := InitCsv(dir); err != nil {
		slog.Error("error initializing report", "error", err)
		return err
	}
	report, err := os.OpenFile(path.Join(dir, "currentWeather.csv"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("file not found, creating empty report file")
		if err = InitCsv(dir); err != nil {
			return err
		}
		report, _ = os.OpenFile(path.Join(dir, "currentWeather.csv"), os.O_APPEND|os.O_WRONLY, 0644)
	}
	defer func() {
		if err := report.Close(); err != nil {
			slog.Error("error closing report stream")
		}
	}()
	writer := csv.NewWriter(report)
	defer writer.Flush()
	reportData = append(reportData, d)
	var chanceOfPrecip string
	if d.ChanceOfPrecip {
		chanceOfPrecip = "true"
	} else {
		chanceOfPrecip = "false"
	}
	for _, data := range reportData {
		row := []string{
			data.Timestamp,
			strconv.FormatFloat(data.Temperature, 'f', 2, 64),
			strconv.FormatFloat(data.Humidity, 'f', 2, 64),
			chanceOfPrecip,
			time.Now().UTC().Format("2006-01-02T15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	slog.Info("successful wrote report", "reportPath", path.Join(dir, "currentWeather.csv"))
	return nil
}
