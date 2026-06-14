package main

import (
	"io"
	"log/slog"
	"os"
	"path"
)

type LogConfig struct {
	LogLevel slog.Level
}

func ConvertCelciusToFahrenheit(temp float64) float64 {
	slog.Debug("converting temp to freedom units", "tempC", temp)
	converted := (temp * 9 / 5) + 32
	return converted
}

func SetLogger(options Environment) (*slog.Logger, error) {
	if err := os.MkdirAll(options.LogDirectory, 0755); err != nil {
		slog.Info("error creating directory")
	}
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
