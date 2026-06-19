package nws

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/espm1000/observing-the-weather/pkg/client"
)

func (n NWSConfig) GetForecastData(h *client.HttpClientConfig) (*ForecastWeatherData, error) {
	var result Forecast

	slog.Info("getting forecast data")
	url := BaseURL + "/gridpoints/" + n.ForecastOffice + "/" + n.GridX + "," + n.GridY + "/forecast"
	resp, err := h.CallGet(url)
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
		slog.Error("non-200 response from upstream", "status_code", resp.StatusCode)
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
