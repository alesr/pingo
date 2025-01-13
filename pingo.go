package pingo

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const defaultURL = "http://www.google.com"

type Options struct {
	frequency  time.Duration
	httpClient *http.Client
	url        string
}

type Option func(*Options)

func WithFrequency(d time.Duration) Option {
	return func(o *Options) {
		o.frequency = d
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(o *Options) {
		o.httpClient = client
	}
}

func WithURL(url string) Option {
	return func(o *Options) {
		o.url = url
	}
}

func StartPinger(ctx context.Context, opts ...Option) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(ctx)

	options := &Options{
		frequency: 45 * time.Second,
		url:       defaultURL,
		httpClient: &http.Client{
			Timeout: 8 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       1,
				MaxConnsPerHost:    1,
				DisableCompression: true,
				DisableKeepAlives:  false,
				IdleConnTimeout:    30 * time.Second,
			},
		},
	}

	for _, opt := range opts {
		opt(options)
	}

	req, err := http.NewRequest(
		http.MethodGet,
		options.url,
		nil,
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	go func() {
		ticker := time.NewTicker(options.frequency)
		defer ticker.Stop()

		ping(options.httpClient, req)

		for {
			select {
			case <-ticker.C:
				ping(options.httpClient, req)
			case <-ctx.Done():
				return
			}
		}
	}()
	return cancel, nil
}

func ping(client *http.Client, req *http.Request) {
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}
