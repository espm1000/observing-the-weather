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
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error making get call", "error", err)
		return nil, err
	}
	return resp, nil
}
