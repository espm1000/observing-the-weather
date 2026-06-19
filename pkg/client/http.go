package client

import (
	"log/slog"
	"net/http"
	"time"
)

type HttpClientConfig struct {
	UserAgent string
	Timeout   time.Duration
	client    http.Client
}

func New(timeout time.Duration, userAgent string) *HttpClientConfig {
	client := http.Client{
		Timeout: timeout,
	}
	return &HttpClientConfig{
		UserAgent: userAgent,
		Timeout:   timeout,
		client:    client,
	}
}

func (h *HttpClientConfig) CallGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("error initiating get request", "error", err)
		return nil, err
	}
	req.Header.Add("User-Agent", h.UserAgent)
	slog.Debug("setting http headers", "headers", req.Header)
	resp, err := h.client.Do(req)
	if err != nil {
		slog.Error("error making get call", "error", err)
		return nil, err
	}
	return resp, nil
}
