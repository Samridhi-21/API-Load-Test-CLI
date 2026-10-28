package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPExecutorSuccess(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"hello"}`))
		}),
	)
	defer server.Close()

	executor := NewHTTPExecutor()

	result := executor.Do(
		context.Background(),
		server.URL,
	)

	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}

	if result.StatusCode != http.StatusOK {
		t.Fatalf(
			"got status %d, want 200",
			result.StatusCode,
		)
	}

	if result.Latency <= 0 {
		t.Fatalf("expected positive latency, got %v", result.Latency)
	}

	if result.Bytes == 0 {
		t.Fatal("expected response body bytes")
	}
}

func TestHTTPExecutorServerError(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.Error(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
		}),
	)
	defer server.Close()

	executor := NewHTTPExecutor()

	result := executor.Do(
		context.Background(),
		server.URL,
	)

	if result.Err != nil {
		t.Fatalf("unexpected transport/read error: %v", result.Err)
	}

	if result.StatusCode != http.StatusInternalServerError {
		t.Fatalf(
			"got status %d, want 500",
			result.StatusCode,
		)
	}
}

func TestHTTPExecutorCancellation(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			select {
			case <-r.Context().Done():
			case <-time.After(2 * time.Second):
			}
		}),
	)
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	executor := NewHTTPExecutor()

	result := executor.Do(ctx, server.URL)

	if result.Err == nil {
		t.Fatal("expected an error for canceled request")
	}
}
