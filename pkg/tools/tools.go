package tools

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"strconv"

	"github.com/caarlos0/env"
	"github.com/espm1000/observing-the-weather/pkg/report"
)

type Environment struct {
	ObservationStationId string `env:"WEATHER_OBSERVATION_STATION_ID" envDefault:"KSTP"`
	ForecastStationId    string `env:"WEATHER_FORECAST_STATION_ID" envDefault:"MPX"`
	LogDirectory         string `env:"WEATHER_LOG_DIRECTORY" envDefault:"logs"`
	LogOutput            string `env:"WEATHER_LOG_FILE" envDefault:"weatherlog.json"`
	PrintToConsole       string `env:"WEATHER_LOG_CONSOLE" envDefault:"false"`
	GridX                string `env:"WEATHER_GRID_X" envDefault:"102"`
	GridY                string `env:"WEATHER_GRID_Y" envDefault:"84"`
	LogLevel             slog.Level
	Logger               *slog.Logger
}

type LogConfig struct {
	LogLevel slog.Level
}

func ConvertCelsiusToFahrenheit(temp float64) (float64, error) {
	slog.Debug("converting temp to fahrenheit", "temp_celsius", temp)
	converted := (temp * 9 / 5) + 32
	return converted, nil
}

func SetEnvironment(e *Environment) error {
	if err := env.Parse(e); err != nil {
		return err
	}
	return nil
}

func SetLogger(options Environment) (*slog.Logger, error) {
	if err := os.MkdirAll(options.LogDirectory, 0755); err != nil {
		slog.Error("error creating directory", "error", err)
		return nil, err
	}
	file, err := os.OpenFile(path.Join(options.LogDirectory, options.LogOutput), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("error creating report folder", "error", err)
		return nil, err
	}
	multiwriter := io.MultiWriter(os.Stdout, file)
	logHandler := slog.NewJSONHandler(multiwriter, &slog.HandlerOptions{
		Level: options.LogLevel,
	})
	logger := slog.New(logHandler)
	return logger, nil
}

func PrintToConsole(d report.CurrentWeatherData, cfg Environment) {
	fmt.Printf("\nCurrent weather for %v \n", d.Timestamp)
	fmt.Printf("Current Temp: %v F\n", strconv.FormatFloat(d.Temperature, 'f', 2, 64))
	fmt.Printf("Current Windspeed: %v km/h\n", d.Windspeed)
	fmt.Printf("Current Humidity: %v Percent\n", strconv.FormatFloat(d.Humidity, 'f', 2, 64))
	fmt.Printf("Grid Coordinates\nX: %v\nY: %v\n", cfg.GridX, cfg.GridY)
	fmt.Printf("Chance of Precip: %v (not implemented)\n", d.ChanceOfPrecip)
}
