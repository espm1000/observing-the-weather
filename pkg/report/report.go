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

type ReportConfig struct {
	Directory  string `env:"WEATHER_REPORT_DIR" envDefault:"/data"`
	ReportFile string `env:"WEATHER_REPORT_FILE" envDefault:"currentWeather.csv"`
}

type CurrentWeatherData struct {
	Temperature    float64
	Humidity       float64
	Windspeed      any
	ChanceOfPrecip bool
	Timestamp      string
	PrecipLastHour float64
}

func InitCsv(r ReportConfig) error {
	slog.Debug("initializing empty csv report")
	headers := []string{"timestamp", "temperature", "humidity", "precipchance", "polledTimestamp"}
	if _, err := os.Stat(r.Directory); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err = os.Mkdir(r.Directory, 0755); err != nil {
				slog.Error("error creating directory", "directory", r.Directory, "error", err)
				return err
			}
		}
	}
	_, err := os.Stat(path.Join(r.Directory, r.ReportFile))
	if err == nil {
		slog.Debug("report file exists")
		return err
	}
	file, err := os.Create(path.Join(r.Directory, r.ReportFile))
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

func WriteCsv(r ReportConfig, d CurrentWeatherData) error {
	var reportData []CurrentWeatherData
	var chanceOfPrecip string
	if err := InitCsv(r); err != nil {
		slog.Error("error initializing report", "error", err)
		return err
	}
	report, err := os.OpenFile(path.Join(r.Directory, r.ReportFile), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("error writing report", "error", err)
		return err
	}
	defer func() {
		if err := report.Close(); err != nil {
			slog.Error("error closing report stream")
		}
	}()

	writer := csv.NewWriter(report)
	defer writer.Flush()

	reportData = append(reportData, d)
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
	slog.Info("successfully wrote report", "reportPath", path.Join(r.Directory, r.ReportFile))
	return nil
}
