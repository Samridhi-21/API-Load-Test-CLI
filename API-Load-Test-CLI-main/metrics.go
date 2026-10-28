package main

import (
	"sync"
	"time"

	"github.com/caio/go-tdigest/v5"
)

// Metrics safely aggregates request outcomes.
type Metrics struct {
	mu sync.Mutex

	totalRequests      int64
	successfulRequests int64
	httpErrors         int64
	transportErrors    int64
	totalBytes         int64
	totalLatency       time.Duration
	latencies          []time.Duration
	digest             *tdigest.TDigest
	retainLatencies    bool
}

// MetricsSnapshot is a consistent copy of the collected metrics.
type MetricsSnapshot struct {
	TotalRequests      int64
	SuccessfulRequests int64
	HTTPErrors         int64
	TransportErrors    int64
	TotalBytes         int64
	TotalLatency       time.Duration
	Latencies          []time.Duration
}

func newMetrics(retainLatencies bool) *Metrics {
	digest, err := tdigest.New()
	if err != nil {
		panic("failed to initialize t-digest: " + err.Error())
	}

	return &Metrics{
		latencies:       make([]time.Duration, 0),
		digest:          digest,
		retainLatencies: retainLatencies,
	}
}

// Exact samples are retained for the Phase 7 baseline.
func NewMetrics() *Metrics {
	return newMetrics(true)
}

// Streaming mode keeps the digest but does not retain every latency.
func NewStreamingMetrics() *Metrics {
	return newMetrics(false)
}

// Record adds one completed request result.
func (m *Metrics) Record(result RequestResult) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests++
	m.totalLatency += result.Latency

	if m.retainLatencies {
		m.latencies = append(m.latencies, result.Latency)
	}

	_ = m.digest.Add(float64(result.Latency))

	m.totalBytes += result.Bytes

	if result.Err != nil {
		m.transportErrors++
		return
	}

	if result.StatusCode >= 200 && result.StatusCode < 400 {
		m.successfulRequests++
	} else if result.StatusCode >= 400 {
		m.httpErrors++
	}
}

// Snapshot returns a consistent copy of the current metrics.
func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	latenciesCopy := append([]time.Duration(nil), m.latencies...)

	return MetricsSnapshot{
		TotalRequests:      m.totalRequests,
		SuccessfulRequests: m.successfulRequests,
		HTTPErrors:         m.httpErrors,
		TransportErrors:    m.transportErrors,
		TotalBytes:         m.totalBytes,
		TotalLatency:       m.totalLatency,
		Latencies:          latenciesCopy,
	}
}

func (s MetricsSnapshot) AverageLatency() time.Duration {
	if s.TotalRequests == 0 {
		return 0
	}

	return time.Duration(
		int64(s.TotalLatency) / s.TotalRequests,
	)
}

func (s MetricsSnapshot) Throughput(elapsed time.Duration) float64 {
	if elapsed <= 0 {
		return 0
	}

	return float64(s.TotalRequests) / elapsed.Seconds()
}

// Quantile returns a t-digest latency estimate.
// p is a fraction between 0 and 1, such as 0.95 for p95.
func (m *Metrics) Quantile(p float64) (time.Duration, bool) {
	if p < 0 || p > 1 {
		return 0, false
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.digest.Count() == 0 {
		return 0, false
	}

	ns := m.digest.Quantile(p)
	return time.Duration(ns), true
}
