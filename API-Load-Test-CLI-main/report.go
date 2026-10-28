package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Report struct {
	Target             string  `json:"target"`
	Concurrency        int     `json:"concurrency"`
	RateLimit          int     `json:"rate_limit"`
	ElapsedMS          float64 `json:"elapsed_ms"`
	TotalRequests      int64   `json:"total_requests"`
	SuccessfulRequests int64   `json:"successful_requests"`
	HTTPErrors         int64   `json:"http_errors"`
	TransportErrors    int64   `json:"transport_errors"`
	TotalBytes         int64   `json:"total_bytes"`
	ThroughputRPS      float64 `json:"throughput_rps"`
	AverageLatencyMS   float64 `json:"average_latency_ms"`
	P50LatencyMS       float64 `json:"p50_latency_ms"`
	P90LatencyMS       float64 `json:"p90_latency_ms"`
	P95LatencyMS       float64 `json:"p95_latency_ms"`
	P99LatencyMS       float64 `json:"p99_latency_ms"`
}

func BuildReport(
	cfg Config,
	elapsed time.Duration,
	metrics *Metrics,
) Report {
	snapshot := metrics.Snapshot()

	quantileMS := func(p float64) float64 {
		value, ok := metrics.Quantile(p)
		if !ok {
			return 0
		}
		return float64(value) / float64(time.Millisecond)
	}

	return Report{
		Target:             cfg.URL,
		Concurrency:        cfg.Concurrency,
		RateLimit:          cfg.Rate,
		ElapsedMS:          float64(elapsed) / float64(time.Millisecond),
		TotalRequests:      snapshot.TotalRequests,
		SuccessfulRequests: snapshot.SuccessfulRequests,
		HTTPErrors:         snapshot.HTTPErrors,
		TransportErrors:    snapshot.TransportErrors,
		TotalBytes:         snapshot.TotalBytes,
		ThroughputRPS:      snapshot.Throughput(elapsed),
		AverageLatencyMS: float64(snapshot.AverageLatency()) /
			float64(time.Millisecond),
		P50LatencyMS: quantileMS(0.50),
		P90LatencyMS: quantileMS(0.90),
		P95LatencyMS: quantileMS(0.95),
		P99LatencyMS: quantileMS(0.99),
	}
}

func WriteJSONReport(path string, report Report) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(report)
}

func WriteCSVReport(path string, report Report) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"target",
		"concurrency",
		"rate_limit",
		"elapsed_ms",
		"total_requests",
		"successful_requests",
		"http_errors",
		"transport_errors",
		"total_bytes",
		"throughput_rps",
		"average_latency_ms",
		"p50_latency_ms",
		"p90_latency_ms",
		"p95_latency_ms",
		"p99_latency_ms",
	}

	row := []string{
		report.Target,
		strconv.Itoa(report.Concurrency),
		strconv.Itoa(report.RateLimit),
		fmt.Sprintf("%.3f", report.ElapsedMS),
		strconv.FormatInt(report.TotalRequests, 10),
		strconv.FormatInt(report.SuccessfulRequests, 10),
		strconv.FormatInt(report.HTTPErrors, 10),
		strconv.FormatInt(report.TransportErrors, 10),
		strconv.FormatInt(report.TotalBytes, 10),
		fmt.Sprintf("%.3f", report.ThroughputRPS),
		fmt.Sprintf("%.3f", report.AverageLatencyMS),
		fmt.Sprintf("%.3f", report.P50LatencyMS),
		fmt.Sprintf("%.3f", report.P90LatencyMS),
		fmt.Sprintf("%.3f", report.P95LatencyMS),
		fmt.Sprintf("%.3f", report.P99LatencyMS),
	}

	if err := writer.Write(header); err != nil {
		return err
	}
	if err := writer.Write(row); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}
