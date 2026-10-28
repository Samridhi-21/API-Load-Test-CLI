package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRenderDashboard(t *testing.T) {
	metrics := NewStreamingMetrics()

	metrics.Record(RequestResult{
		Latency:    20 * time.Millisecond,
		StatusCode: 200,
	})

	var output bytes.Buffer

	RenderDashboard(
		&output,
		metrics,
		time.Now().Add(-time.Second),
		10*time.Second,
	)

	got := output.String()

	for _, expected := range []string{
		"API LOAD TEST",
		"Elapsed:",
		"Remaining:",
		"REQUESTS",
		"Successful:",
		"Throughput:",
		"LATENCY",
		"p50:",
		"p95:",
		"p99:",
	} {
		if !strings.Contains(got, expected) {
			t.Errorf("dashboard output missing %q", expected)
		}
	}
}
