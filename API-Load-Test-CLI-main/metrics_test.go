package main

import (
	"errors"
	"testing"
	"time"
)

func TestMetricsRecordAndSnapshot(t *testing.T) {
	metrics := NewMetrics()

	metrics.Record(RequestResult{
		Latency:    10 * time.Millisecond,
		StatusCode: 200,
		Bytes:      100,
	})

	metrics.Record(RequestResult{
		Latency:    20 * time.Millisecond,
		StatusCode: 302,
		Bytes:      50,
	})

	metrics.Record(RequestResult{
		Latency:    30 * time.Millisecond,
		StatusCode: 500,
		Bytes:      25,
	})

	metrics.Record(RequestResult{
		Latency: 5 * time.Millisecond,
		Err:     errors.New("connection refused"),
	})

	snapshot := metrics.Snapshot()

	if snapshot.TotalRequests != 4 {
		t.Errorf("total requests = %d, want 4",
			snapshot.TotalRequests)
	}

	if snapshot.SuccessfulRequests != 2 {
		t.Errorf("successful = %d, want 2",
			snapshot.SuccessfulRequests)
	}

	if snapshot.HTTPErrors != 1 {
		t.Errorf("HTTP errors = %d, want 1",
			snapshot.HTTPErrors)
	}

	if snapshot.TransportErrors != 1 {
		t.Errorf("transport errors = %d, want 1",
			snapshot.TransportErrors)
	}

	if snapshot.TotalBytes != 175 {
		t.Errorf("total bytes = %d, want 175",
			snapshot.TotalBytes)
	}

	if len(snapshot.Latencies) != 4 {
		t.Errorf("latency samples = %d, want 4",
			len(snapshot.Latencies))
	}

	wantAverage := 16*time.Millisecond + 250*time.Microsecond
	if snapshot.AverageLatency() != wantAverage {
		t.Errorf("average latency = %v, want %v",
			snapshot.AverageLatency(), wantAverage)
	}
}

func TestMetricsEmptySnapshot(t *testing.T) {
	snapshot := NewMetrics().Snapshot()

	if snapshot.AverageLatency() != 0 {
		t.Errorf("empty average latency should be zero")
	}

	if snapshot.Throughput(time.Second) != 0 {
		t.Errorf("empty throughput should be zero")
	}
}
