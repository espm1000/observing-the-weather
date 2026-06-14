package tools

import (
	"io"
	"log/slog"
	"os"
	"path"
)

type Environment struct {
	ReportOutputDir      string `env:"WEATHER_REPORT_DIR" envDefault:"/data"`
	ObservationStationId string `env:"WEATHER_OBSERVATION_STATION_ID" envDefault:"KSTP"`
	ForecastStationId    string `env:"WEATHER_FORECAST_STATION_ID" envDefault:"MPX"`
	LogDirectory         string `env:"WEATHER_LOG_DIRECTORY" envDefault:"logs"`
	LogOutput            string `env:"WEATHER_LOG_FILE" envDefault:"weatherlog.json"`
	LogLevel             slog.Level
}

type LogConfig struct {
	LogLevel slog.Level
}

func ConvertCelciusToFahrenheit(temp float64) (float64, error) {
	slog.Info("converting temp to freedom units", "tempC", temp)
	converted := (temp * 9 / 5) + 32
	return converted, nil
}

func SetLogger(options Environment) (*slog.Logger, error) {
	if err := os.MkdirAll(options.LogDirectory, 0755); err != nil {
		slog.Info("error creating directory")
	}
	slog.Info("setting log config", "log_directory", options.LogDirectory)
	file, err := os.OpenFile(path.Join(options.LogDirectory, "weatherlog.json"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	multiwriter := io.MultiWriter(os.Stdout, file)
	logHandler := slog.NewJSONHandler(multiwriter, &slog.HandlerOptions{
		Level: options.LogLevel,
	})
	logger := slog.New(logHandler)
	return logger, nil
}
