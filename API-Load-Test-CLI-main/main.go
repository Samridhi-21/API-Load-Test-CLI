package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	cfg, err := ParseConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(2)
	}

	executor := NewHTTPExecutor()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.Duration,
	)
	defer cancel()

	start := time.Now()
	results := RunWorkers(ctx, cfg, executor)

	metrics := NewStreamingMetrics()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for results != nil {
		select {
		case result, ok := <-results:
			if !ok {
				results = nil
				continue
			}

			metrics.Record(result)

		case <-ticker.C:
			RenderDashboard(
				os.Stdout,
				metrics,
				start,
				cfg.Duration,
			)
		}
	}

	// Render the final live snapshot before the final summary.
	RenderDashboard(
		os.Stdout,
		metrics,
		start,
		cfg.Duration,
	)

	elapsed := time.Since(start)

	report := BuildReport(cfg, elapsed, metrics)

	if cfg.JSONOutput != "" {
		if err := WriteJSONReport(cfg.JSONOutput, report); err != nil {
			fmt.Fprintln(os.Stderr, "JSON report error:", err)
		} else {
			fmt.Println("JSON report:", cfg.JSONOutput)
		}
	}

	if cfg.CSVOutput != "" {
		if err := WriteCSVReport(cfg.CSVOutput, report); err != nil {
			fmt.Fprintln(os.Stderr, "CSV report error:", err)
		} else {
			fmt.Println("CSV report:", cfg.CSVOutput)
		}
	}

	snapshot := metrics.Snapshot()

	fmt.Println("\nAPI Load Test Complete")
	fmt.Println("----------------------")
	fmt.Println("Target:             ", cfg.URL)
	fmt.Println("Workers:            ", cfg.Concurrency)
	fmt.Println("Rate limit:         ", cfg.Rate, "requests/sec")
	fmt.Println("Elapsed:            ", elapsed.Round(time.Millisecond))
	fmt.Println("Total requests:     ", snapshot.TotalRequests)
	fmt.Println("Successful:         ", snapshot.SuccessfulRequests)
	fmt.Println("HTTP errors:        ", snapshot.HTTPErrors)
	fmt.Println("Transport errors:   ", snapshot.TransportErrors)
	fmt.Println("Total bytes:        ", snapshot.TotalBytes)
	fmt.Println("Average latency:    ", snapshot.AverageLatency())
	fmt.Printf(
		"Throughput:          %.2f req/s\n",
		snapshot.Throughput(elapsed),
	)

	fmt.Println("\nStreaming latency percentiles")

	for _, p := range []float64{0.50, 0.90, 0.95, 0.99} {
		latency, ok := metrics.Quantile(p)

		if !ok {
			fmt.Printf("p%.0f: no samples\n", p*100)
			continue
		}

		fmt.Printf("p%.0f: %v\n", p*100, latency)
	}
}
