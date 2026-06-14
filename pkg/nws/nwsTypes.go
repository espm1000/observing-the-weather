package nws

import (
	"time"
)

type NWSConfig struct {
	BaseURL        string
	GridX          string
	GridY          string
	ForecastOffice string
	StationID      string
}

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

type QuantitativeValue struct {
	Value          float64 `json:"value"`
	MaxValue       float64 `json:"maxValue"`
	MinValue       float64 `json:"minValue"`
	UnitCode       string  `json:"unitCode"`
	QualityControl string  `json:"qualityControl"`
}

type ObservationQuantitativeValue struct {
	Value          float64 `json:"value"`
	MaxValue       float64 `json:"maxValue"`
	MinValue       float64 `json:"minValue"`
	UnitCode       string  `json:"unitCode"`
	QualityControl string  `json:"qualityControl"`
}

type PresentWeather struct {
	Intensity  string `json:"intensity"`
	Modifier   string `json:"modifier"`
	Weather    string `json:"weather"`
	RawString  string `json:"rawString"`
	InVicinity bool   `json:"inVicinity"`
}

type CloudLayer struct {
	Base   QuantitativeValue `json:"base"`
	Amount string            `json:"amount"`
}

type ObservationProperties struct {
	Context                   any                          `json:"@context"`
	Geometry                  string                       `json:"geometry"`
	ID                        string                       `json:"@id"`
	Type                      string                       `json:"@type"`
	Elevation                 QuantitativeValue            `json:"elevation"`
	Station                   string                       `json:"station"`
	StationID                 string                       `json:"stationId"`
	StationName               string                       `json:"stationName"`
	Timestamp                 string                       `json:"timestamp"`
	RawMessage                string                       `json:"rawMessage"`
	TextDescription           string                       `json:"textDescription"`
	PresentWeather            []PresentWeather             `json:"presentWeather"`
	Temperature               ObservationQuantitativeValue `json:"temperature"`
	Dewpoint                  ObservationQuantitativeValue `json:"dewpoint"`
	WindSpeed                 ObservationQuantitativeValue `json:"windSpeed"`
	Visibility                ObservationQuantitativeValue `json:"visibility"`
	MaxTemperatureLast24Hours ObservationQuantitativeValue `json:"maxTemperatureLast24Hours"`
	MinTemperatureLast24Hours ObservationQuantitativeValue `json:"minTemperatureLast24Hours"`
	PrecipitationLastHour     ObservationQuantitativeValue `json:"precipitationLastHour"`
	PrecipitationLast3Hours   ObservationQuantitativeValue `json:"precipitationLast3Hours"`
	PrecipitationLast6Hours   ObservationQuantitativeValue `json:"precipitationLast6Hours"`
	RelativeHumidity          ObservationQuantitativeValue `json:"relativeHumidity"`
	WindChill                 ObservationQuantitativeValue `json:"windChill"`
	HeatIndex                 ObservationQuantitativeValue `json:"heatIndex"`
	CloudLayers               []CloudLayer                 `json:"cloudLayers"`
}

type ObservationGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
	Bbox        []float64 `json:"bbox"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates any       `json:"coordinates"`
	BBox        []float64 `json:"bbox"`
}

type Observation struct {
	Context    any                   `json:"@context"`
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Geometry   Geometry              `json:"geometry"`
	Properties ObservationProperties `json:"properties"`
}
