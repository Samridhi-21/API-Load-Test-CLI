package main

import (
	"reflect"
	"testing"
	"time"
)

func TestPercentileNearestRank(t *testing.T) {
	latencies := []time.Duration{
		100 * time.Millisecond,
		10 * time.Millisecond,
		40 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
	}

	tests := []struct {
		name string
		p    float64
		want time.Duration
	}{
		{"p50", 50, 30 * time.Millisecond},
		{"p90", 90, 100 * time.Millisecond},
		{"p95", 95, 100 * time.Millisecond},
		{"p99", 99, 100 * time.Millisecond},
		{"p100", 100, 100 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Percentile(latencies, tt.p)

			if !ok {
				t.Fatal("expected a valid percentile")
			}

			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPercentileDoesNotModifyInput(t *testing.T) {
	latencies := []time.Duration{
		3 * time.Millisecond,
		1 * time.Millisecond,
		2 * time.Millisecond,
	}

	original := append([]time.Duration(nil), latencies...)

	_, ok := Percentile(latencies, 50)
	if !ok {
		t.Fatal("expected a valid percentile")
	}

	if !reflect.DeepEqual(latencies, original) {
		t.Fatal("Percentile modified the original slice")
	}
}

func TestPercentileInvalidInput(t *testing.T) {
	_, ok := Percentile(nil, 50)
	if ok {
		t.Fatal("empty samples should be rejected")
	}

	for _, p := range []float64{0, -1, 101} {
		_, ok := Percentile(
			[]time.Duration{time.Millisecond},
			p,
		)
		if ok {
			t.Errorf("percentile %v should be rejected", p)
		}
	}
}
