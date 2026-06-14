package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

func GetHistoricalObservations(c NWSConfig) (*ObservationCollection, error) {
	slog.Info("getting historical observations", "startDate", c.StartDate, "endDate", c.EndDate)
	var collection ObservationCollection
	params := url.Values{}
	params.Add("end", c.EndDate)
	params.Add("start", c.StartDate)
	uri := fmt.Sprintf(HistoricalURI, c.StationID)
	url, err := url.Parse(BaseURL + uri)
	if err != nil {
		slog.Error("error parsing url")
	}
	url.RawQuery = params.Encode()

	nextURL := url
	nextURL_s := nextURL.String()
	obs := make(map[string]HistoricalObvs)
	for nextURL_s != "" {
		slog.Info("next token received")
		resp, err := http.Get(nextURL_s)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				slog.Error("error closing stream")
			}
		}()
		if err != nil {
			slog.Error("error calling http", "error", err)
			return nil, err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("error reading body", "error", err)
			return nil, err
		}
		var collection ObservationCollection
		if err := json.Unmarshal(body, &collection); err != nil {
			slog.Error("error jsoning", "error", err)
			return nil, err
		}
		ParseObservations(collection)

		for _, o := range collection.Features {
			obs[o.Properties.Timestamp] = HistoricalObvs{
				Temperature: o.Properties.Temperature.Value,
				Humidity:    o.Properties.RelativeHumidity.Value,
			}
		}
		nextURL_s = collection.Pagination.Next
	}
	// ParseObservations(collection)
	// req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	// if err != nil {
	// 	return nil, err
	// }
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	slog.Error("error calling url")
	// 	return nil, err
	// }
	// defer resp.Body.Close()
	// json.NewDecoder(resp.Body).Decode(&collection)
	// token_url := collection.Pagination.Next
	// parsedUrl, err := url.Parse(token_url)
	// if err != nil {
	// 	slog.Error("error parsing url", "error", err)
	// }
	// cursor := parsedUrl.Query().Get("cursor")
	// if len(cursor) > 0 {
	// 	fmt.Println("more")
	// }
	// // fmt.Println(cursor)
	// if err := ParseObservations(collection); err != nil {
	// 	slog.Error("error parsing observations")
	// 	return nil, err
	// }

	return &collection, err
}
