package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
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
	ReportOutputDir string
}

func ReadEnvValues() Environment {
	return Environment{
		ReportOutputDir: os.Getenv("WEATHER_REPORT_DIR"),
	}
}

func (n NWSConfig) GetCurrentData() (*CurrentWeatherData, error) {
	var currentData Observation
	slog.Info("getting current weather data", "station", n.StationID)
	resp, err := http.Get(n.BaseURL + "/stations/" + n.StationID + "/observations/latest")
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
	temp_f, err := ConvertCelciusToFahrenheit(currentData.Properties.Temperature.Value)
	if err != nil {
		return nil, err
	}

	return &CurrentWeatherData{
		Temperature: temp_f,
		Humidity:    currentData.Properties.RelativeHumidity.Value,
		Windspeed:   currentData.Properties.WindSpeed.Value,
		Timestamp:   currentData.Properties.Timestamp,
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

func InitCsv(dir string) error {
	headers := []string{"timestamp", "temperature", "humidity"}
	_, err := os.Stat(path.Join(dir, "currentWeather.csv"))
	if err == nil {
		fmt.Println("report file exists")
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
	for _, data := range reportData {
		row := []string{
			data.Timestamp,
			strconv.FormatFloat(data.Temperature, 'f', 2, 64),
			strconv.FormatFloat(data.Humidity, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	slog.Info("successful wrote report", "reportPath", dir+"currentWeather.csv")
	return nil
}

func PrintToConsole(d CurrentWeatherData) {
	fmt.Printf("\nCurrent weather for %v \n", d.Timestamp)
	fmt.Printf("Current Temp: %v F\n", d.Temperature)
	fmt.Printf("Current Windspeed: %v km/h\n", d.Windspeed)
	fmt.Printf("Current Humidity: %v Percent\n", strconv.FormatFloat(d.Humidity, 'f', 2, 64))
}

func main() {
	var reportDir string
	env := ReadEnvValues()
	if env != (Environment{}) {
		slog.Info("setting options per environment variables")
		reportDir = env.ReportOutputDir
	} else {
		slog.Info("using default values")
		reportDir = "/data"
	}
	nws := NWSConfig{
		BaseURL:        "https://api.weather.gov",
		GridX:          "102",
		GridY:          "84",
		ForecastOffice: "MPX",  // Minneapolis
		StationID:      "KSTP", // St. Paul
	}
	CurrentWeather, err := nws.GetCurrentData()
	if err != nil {
		slog.Error("error", "error", err)
		panic(err)
	}
	if err := WriteCsv(reportDir, *CurrentWeather); err != nil {
		log.Fatal(err)
	}

	PrintToConsole(*CurrentWeather)
}
