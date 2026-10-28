package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPoolCompletesRequests(t *testing.T) {
	var requests atomic.Int64

	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			requests.Add(1)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}),
	)
	defer server.Close()

	cfg := Config{
		URL:         server.URL,
		Concurrency: 5,
		Duration:    200 * time.Millisecond,
		Rate:        0,
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.Duration,
	)
	defer cancel()

	results := RunWorkers(ctx, cfg, NewHTTPExecutor())

	var completed int
	for range results {
		completed++
	}

	if completed == 0 {
		t.Fatal("expected completed request attempts")
	}

	if requests.Load() == 0 {
		t.Fatal("server received no requests")
	}

	t.Logf(
		"Completed results: %d; server requests: %d",
		completed,
		requests.Load(),
	)
}
