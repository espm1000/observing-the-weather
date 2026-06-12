package main

import (
	"encoding/json"
	"fmt"
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
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("error calling url")
		return nil, err
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&collection)
	token_url := collection.Pagination.Next
	parsedUrl, err := url.Parse(token_url)
	if err != nil {
		slog.Error("error parsing url", "error", err)
	}
	cursor := parsedUrl.Query().Get("cursor")
	if len(cursor) > 0 {
		fmt.Println("more")
	}
	// fmt.Println(cursor)
	// if err := ParseObservations(collection); err != nil {
	// 	slog.Error("error parsing observations")
	// 	return nil, err
	// }

	return &collection, err
}
