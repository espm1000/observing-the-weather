package nws

import (
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
