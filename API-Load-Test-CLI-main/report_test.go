package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteJSONReport(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "report.json")

	report := Report{
		Target:        "http://localhost:8080/",
		TotalRequests: 42,
		ThroughputRPS: 10.5,
	}

	if err := WriteJSONReport(path, report); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var got Report
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}

	if got.TotalRequests != 42 {
		t.Fatalf("total requests = %d, want 42", got.TotalRequests)
	}
}

func TestWriteCSVReport(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.csv")

	report := Report{
		Target:        "http://localhost:8080/",
		TotalRequests: 42,
	}

	if err := WriteCSVReport(path, report); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 2 {
		t.Fatalf("CSV rows = %d, want header and one data row", len(rows))
	}

	if rows[1][4] != "42" {
		t.Errorf("CSV total requests = %q, want 42", rows[1][4])
	}
}

func TestBuildReport(t *testing.T) {
	cfg := Config{
		URL:         "http://localhost:8080/",
		Concurrency: 4,
		Rate:        10,
	}

	metrics := NewStreamingMetrics()
	metrics.Record(RequestResult{
		Latency:    10 * time.Millisecond,
		StatusCode: 200,
	})

	report := BuildReport(cfg, time.Second, metrics)

	if report.TotalRequests != 1 {
		t.Errorf("total requests = %d, want 1", report.TotalRequests)
	}

	if report.SuccessfulRequests != 1 {
		t.Errorf("successful requests = %d, want 1",
			report.SuccessfulRequests)
	}

	if report.ThroughputRPS != 1 {
		t.Errorf("throughput = %v, want 1", report.ThroughputRPS)
	}
}
