package main

import (
	"testing"
	"time"
)

func TestStreamingMetricsTracksSamples(t *testing.T) {
	metrics := NewStreamingMetrics()

	for i := 1; i <= 1000; i++ {
		metrics.Record(RequestResult{
			Latency:    time.Duration(i) * time.Millisecond,
			StatusCode: 200,
		})
	}

	snapshot := metrics.Snapshot()

	if snapshot.TotalRequests != 1000 {
		t.Fatalf("total requests = %d, want 1000",
			snapshot.TotalRequests)
	}

	// Streaming mode should not retain all individual samples.
	if len(snapshot.Latencies) != 0 {
		t.Fatalf("retained %d individual latencies, want 0",
			len(snapshot.Latencies))
	}

	// The digest should provide a plausible median estimate.
	p50, ok := metrics.Quantile(0.50)
	if !ok {
		t.Fatal("expected a p50 estimate")
	}

	if p50 < 300*time.Millisecond ||
		p50 > 700*time.Millisecond {
		t.Errorf("unexpected p50 estimate: %v", p50)
	}

	p99, ok := metrics.Quantile(0.99)
	if !ok {
		t.Fatal("expected a p99 estimate")
	}

	if p99 < 800*time.Millisecond ||
		p99 > 1100*time.Millisecond {
		t.Errorf("unexpected p99 estimate: %v", p99)
	}
}

func TestStreamingQuantileValidation(t *testing.T) {
	metrics := NewStreamingMetrics()

	if _, ok := metrics.Quantile(0.95); ok {
		t.Fatal("empty digest should not return a quantile")
	}

	metrics.Record(RequestResult{
		Latency:    10 * time.Millisecond,
		StatusCode: 200,
	})

	for _, p := range []float64{-0.1, 1.1} {
		if _, ok := metrics.Quantile(p); ok {
			t.Errorf("quantile %v should be rejected", p)
		}
	}
}
