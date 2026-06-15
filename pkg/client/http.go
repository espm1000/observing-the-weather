package client

import (
	"log/slog"
	"net/http"
	"time"
)

func setClient() http.Client {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	return client
}

func CallGet(url string) (*http.Response, error) {
	client := setClient()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("error initiating get request", "error", err)
		return nil, err
	}
	req.Header.Add("User-Agent", "weather@esp.m1k@gmail.com")
	slog.Debug("setting http headers", "headers", req.Header)
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error making get call", "error", err)
		return nil, err
	}
	return resp, nil
}
