package main

import (
	"fmt"
	"io"
	"time"
)

func RenderDashboard(
	w io.Writer,
	metrics *Metrics,
	start time.Time,
	duration time.Duration,
) {
	elapsed := time.Since(start)
	if elapsed < 0 {
		elapsed = 0
	}

	remaining := duration - elapsed
	if remaining < 0 {
		remaining = 0
	}

	progress := float64(elapsed) / float64(duration)
	if progress > 1 {
		progress = 1
	}

	const barWidth = 30
	filled := int(progress * barWidth)

	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "="
		} else {
			bar += "-"
		}
	}

	snapshot := metrics.Snapshot()

	p50, hasP50 := metrics.Quantile(0.50)
	p95, hasP95 := metrics.Quantile(0.95)
	p99, hasP99 := metrics.Quantile(0.99)

	// Clear the previous dashboard frame.
	fmt.Fprint(w, "\033[H\033[2J")

	fmt.Println() // Keep output aligned in ordinary terminal use.
	fmt.Fprintln(w, "API LOAD TEST — LIVE DASHBOARD")
	fmt.Fprintln(w, "================================")
	fmt.Fprintf(w, "Elapsed:       %v\n", elapsed.Round(time.Millisecond))
	fmt.Fprintf(w, "Remaining:     %v\n", remaining.Round(time.Millisecond))
	fmt.Fprintf(w, "[%s] %3.0f%%\n", bar, progress*100)

	fmt.Fprintln(w, "\nREQUESTS")
	fmt.Fprintf(w, "Total:         %d\n", snapshot.TotalRequests)
	fmt.Fprintf(w, "Successful:    %d\n", snapshot.SuccessfulRequests)
	fmt.Fprintf(w, "HTTP errors:   %d\n", snapshot.HTTPErrors)
	fmt.Fprintf(w, "Transport err: %d\n", snapshot.TransportErrors)
	fmt.Fprintf(w, "Throughput:    %.2f req/s\n",
		snapshot.Throughput(elapsed))

	fmt.Fprintln(w, "\nLATENCY")
	fmt.Fprintf(w, "Average:       %v\n", snapshot.AverageLatency())
	printQuantile(w, "p50", p50, hasP50)
	printQuantile(w, "p95", p95, hasP95)
	printQuantile(w, "p99", p99, hasP99)

	fmt.Fprintln(w)
}

func printQuantile(
	w io.Writer,
	name string,
	value time.Duration,
	ok bool,
) {
	if !ok {
		fmt.Fprintf(w, "%-14s%s\n", name+":", "n/a")
		return
	}

	fmt.Fprintf(w, "%-14s%v\n", name+":", value)
}
