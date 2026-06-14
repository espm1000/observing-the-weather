package forecast

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type ForecastWeatherData struct {
	Temperature    any
	Humidity       float64
	Windspeed      any
	ChanceOfPrecip bool
	PrecipPercent  any
	Timestamp      time.Time
}

type Forecast struct {
	Context    any           `json:"@context"`
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Geometry   Geometry      `json:"geometry"`
	Properties ForecastProps `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates any       `json:"coordinates"`
	BBox        []float64 `json:"bbox"`
}

type QuantitativeValue struct {
	Value          float64 `json:"value"`
	MaxValue       float64 `json:"maxValue"`
	MinValue       float64 `json:"minValue"`
	UnitCode       string  `json:"unitCode"`
	QualityControl string  `json:"qualityControl"`
}

type ForecastPeriod struct {
	Number                     int               `json:"number"`
	Name                       string            `json:"name"`
	StartTime                  time.Time         `json:"startTime"`
	EndTime                    time.Time         `json:"endTime"`
	IsDaytime                  bool              `json:"isDaytime"`
	Temperature                any               `json:"temperature"`
	TemperatureTrend           string            `json:"temperatureTrend"`
	ProbabilityOfPrecipitation QuantitativeValue `json:"probabilityOfPrecipitation"`
	WindSpeed                  any               `json:"windSpeed"`
}

type ForecastProps struct {
	Context           any               `json:"@context"`
	Geometry          string            `json:"geometry"`
	Units             string            `json:"units"`
	ForecastGenerator string            `json:"forecastGenerator"`
	GeneratedAt       time.Time         `json:"generatedAt"`
	UpdateTime        time.Time         `json:"updateTime"`
	ValidTimes        string            `json:"validTimes"`
	Elevation         QuantitativeValue `json:"elevation"`
	Periods           []ForecastPeriod  `json:"periods"`
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
