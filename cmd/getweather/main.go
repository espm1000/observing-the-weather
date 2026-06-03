package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
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

func InitCsv() error {
	headers := []string{"timestamp", "temperature", "humidity"}
	_, err := os.Stat("currentWeather.csv")
	if err == nil {
		fmt.Println("report file exists")
		return err
	}
	file, err := os.Create("currentWeather.csv")
	if err != nil {
		slog.Error("error creating current report file", "error", err)
		return err
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(headers)
	return err
}

func WriteCsv(d CurrentWeatherData) error {
	var reportData []CurrentWeatherData
	// headers := []string{"timestamp", "temperature", "humidity"}
	// report, err := os.Create("currentWeather.csv")
	// if err != nil {
	// 	slog.Error("failed to create file", "error", err)
	// 	return err
	// }
	// defer report.Close()
	report, err := os.OpenFile("currentWeather.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("file not found, creating empty report file")
		InitCsv()
		report, _ = os.OpenFile("currentWeather.csv", os.O_APPEND|os.O_WRONLY, 0644)
	}
	defer report.Close()
	writer := csv.NewWriter(report)
	defer writer.Flush()
	// if err := writer.Write(headers); err != nil {
	// 	slog.Error("error writing headers", "error", err)
	// 	return err
	// }
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
	return nil
}

func PrintToConsole(d CurrentWeatherData) error {
	tempF, err := ConvertCelciusToFahrenheit(d.Temperature)
	if err != nil {
		return err
	}
	fmt.Printf("\nCurrent weather for %v \n", d.Timestamp)
	fmt.Printf("Current Temp: %v F\n", tempF)
	fmt.Printf("Current Windspeed: %v km/h\n", d.Windspeed)
	fmt.Printf("Current Humidity: %v Percent\n", strconv.FormatFloat(d.Humidity, 'f', 2, 64))
	return err
}

func main() {
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
	if err := WriteCsv(*CurrentWeather); err != nil {
		log.Fatal(err)
	}

	// if err := PrintToConsole(*CurrentWeather); err != nil {
	// 	slog.Error("error printing to console", "error", err)
	// }
}
