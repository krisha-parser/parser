package scraper

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const (
	defaultMaxRetries = 3
	defaultRetryDelay = time.Second
)

type loadOptions struct {
	maxRetries int
	retryDelay time.Duration
	headers    http.Header
}

type Option func(*loadOptions)

func withMaxRetries(n int) Option {
	return func(o *loadOptions) {
		o.maxRetries = n
	}
}

func withRetryDelay(d time.Duration) Option {
	return func(o *loadOptions) {
		o.retryDelay = d
	}
}

func withHeaders(h http.Header) Option {
	return func(o *loadOptions) {
		o.headers = h
	}
}

func load(client *http.Client, url string, opts ...Option) (string, error) {
	slog.Debug("loading url: " + url)
	options := &loadOptions{
		maxRetries: defaultMaxRetries,
		retryDelay: defaultRetryDelay,
		headers:    make(http.Header),
	}
	for _, opt := range opts {
		opt(options)
	}

	var lastErr error
	for attempt := range options.maxRetries {
		if attempt > 0 {
			time.Sleep(options.retryDelay)
		}

		body, err := doRequest(client, url, options)
		if err == nil {
			return body, nil
		}
		lastErr = err
	}

	return "", fmt.Errorf("all %d attempts failed, last error: %w", options.maxRetries, lastErr)
}

func doRequest(client *http.Client, url string, opts *loadOptions) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	for key, values := range opts.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			panic(err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return "", errors.New(response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
