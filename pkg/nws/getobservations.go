package nws

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/espm1000/observing-the-weather/pkg/client"
	"github.com/espm1000/observing-the-weather/pkg/report"
	"github.com/espm1000/observing-the-weather/pkg/tools"
)

func (n NWSConfig) GetCurrentData(h *client.HttpClientConfig) (*report.CurrentWeatherData, error) {
	var currentData Observation
	slog.Info("getting current weather data", "observationStation", n.StationID, "forecastOffice", n.ForecastOffice)
	url := BaseURL + "/stations/" + n.StationID + "/observations/latest"
	resp, err := h.CallGet(url)
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
	temp_f, err := tools.ConvertCelciusToFahrenheit(currentData.Properties.Temperature.Value)
	if err != nil {
		return nil, err
	}

	return &report.CurrentWeatherData{
		Temperature:    temp_f,
		Humidity:       currentData.Properties.RelativeHumidity.Value,
		Windspeed:      currentData.Properties.WindSpeed.Value,
		Timestamp:      currentData.Properties.Timestamp,
		ChanceOfPrecip: false,
	}, nil

}
