package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strconv"
	"time"
)

func InitCsv(dir string) error {
	slog.Info("initializing empty csv report")
	headers := []string{"timestamp", "temperature", "humidity", "precipchance", "polledTimestamp"}
	_, err := os.Stat(path.Join(dir, "currentWeather.csv"))
	if err == nil {
		slog.Info("report file exists")
		return err
	}
	file, err := os.Create(path.Join(dir, "currentWeather.csv"))
	if err != nil {
		slog.Error("error creating current report file", "error", err)
		return err
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write(headers); err != nil {
		return err
	}
	return err
}

func WriteCsv(dir string, d CurrentWeatherData) error {
	var reportData []CurrentWeatherData
	_, err := os.Stat(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(dir, 0755); err != nil {
				slog.Error("error creating directory", "directory", dir)
				return err
			}
		}
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
			data.Timestamp, // timestamp
			strconv.FormatFloat(data.Temperature, 'f', 2, 64), // temperature
			strconv.FormatFloat(data.Humidity, 'f', 2, 64),    // humidity
			chanceOfPrecip, // chanceofprecip
			time.Now().UTC().Format("2006-01-02T15:04:05"), // polledTimestamp
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	slog.Info("successful wrote report", "reportPath", dir+"currentWeather.csv")
	return nil
}

type HistoricalObvs struct {
	Temperature float64
	Humidity    float64
}

type HistoricalCollection struct {
	Timestamp string
	Data      HistoricalObvs
}

func ParseObservations(o ObservationCollection) error {
	parsedCollection := make(map[string]HistoricalObvs)
	for _, data := range o.Features {
		parsedCollection[data.Properties.Timestamp] = HistoricalObvs{
			Temperature: ConvertCelciusToFahrenheit(data.Properties.Temperature.Value),
			Humidity:    data.Properties.RelativeHumidity.Value,
		}
	}
	slog.Info("sending to write csv")
	if err := WriteObservationsReport("data", parsedCollection); err != nil {
		slog.Error("error writing report", "error", err)
	}
	return nil
}

func WriteObservationsReport(dir string, observations map[string]HistoricalObvs) error {
	if err := InitCsv(dir); err != nil {
		slog.Error("error initializing csv", "error", err)
	}
	file, err := os.OpenFile(path.Join(dir, "currentWeather.csv"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("error opening file", "error", err)
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for ts, data := range observations {
		row := []string{
			ts,
			fmt.Sprintf("%.2f", data.Temperature),
			fmt.Sprintf("%.2f", data.Humidity),
			"false",
			time.Now().UTC().Format("2006-01-02T15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			slog.Error("error writing rows", "error", err)
			return err
		}
	}
	return writer.Error()
}
