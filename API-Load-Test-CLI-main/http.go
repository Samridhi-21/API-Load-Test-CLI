package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

type HTTPExecutor struct {
	client *http.Client
}

type RequestResult struct {
	Latency    time.Duration
	StatusCode int
	Bytes      int64
	Err        error
}

func NewHTTPExecutor() *HTTPExecutor {
	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 500,
		MaxConnsPerHost:     0,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}

	return &HTTPExecutor{client: client}
}

func (e *HTTPExecutor) Do(
	ctx context.Context,
	targetURL string,
) RequestResult {
	start := time.Now()

	result := RequestResult{}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		targetURL,
		nil,
	)
	if err != nil {
		result.Latency = time.Since(start)
		result.Err = err
		return result
	}

	resp, err := e.client.Do(req)
	if err != nil {
		result.Latency = time.Since(start)
		result.Err = err
		return result
	}

	result.StatusCode = resp.StatusCode

	// Read the response body so the connection can be reused.
	result.Bytes, result.Err = io.Copy(io.Discard, resp.Body)

	closeErr := resp.Body.Close()
	if result.Err == nil {
		result.Err = closeErr
	}

	result.Latency = time.Since(start)

	return result
}
