package main

import "log/slog"

func ConvertCelciusToFahrenheit(temp float64) (float64, error) {
	slog.Info("converting temp to freedom units", "tempC", temp)
	converted := (temp * 9 / 5) + 32
	return converted, nil
}
