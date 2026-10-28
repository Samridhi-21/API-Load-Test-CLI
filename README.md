# APIStress — Lightweight API Load Testing Tool

A lightweight command-line tool built in Go for testing API performance under concurrent HTTP traffic. The tool measures response times, throughput, and errors while a test is running, and generates reports for further analysis.

## Overview

APIs need to remain responsive when multiple clients send requests simultaneously. This project aims to provide a simple way to measure API performance under load without relying on a GUI-based testing environment.

APIStress allows users to configure the target endpoint, concurrency level, test duration, and request rate. It collects performance metrics, estimates response-time percentiles, displays live results in the terminal, and exports the final measurements.

## Objectives

- Generate concurrent HTTP requests using Go.
- Control request traffic using a configurable rate limiter.
- Measure response latency, throughput, and HTTP errors.
- Calculate p50, p95, and p99 latency percentiles.
- Use a streaming t-digest estimator to reduce the need to retain every latency sample.
- Display performance metrics in a live terminal dashboard.
- Export test results in JSON and CSV formats.

## Key Features

- **Command-line interface:** Configure and run tests from the terminal.
- **Concurrent request generation:** Use a worker pool to manage multiple requests.
- **Rate limiting:** Control the target request-start rate.
- **Performance metrics:** Measure throughput, response times, and error rates.
- **Percentile analysis:** Examine typical and slower API responses using p50, p95, and p99.
- **Live monitoring:** View performance metrics during test execution.
- **Report generation:** Save results in structured JSON and CSV formats.


## Technology Stack

| Technology | Purpose |
|---|---|
| Go | Main programming language |
| `net/http` | HTTP request execution |
| Goroutines and channels | Concurrent request processing and metrics collection |
| `context` | Cancellation and test lifecycle management |
| Cobra, if used | Command-line argument parsing |
| t-digest library | Streaming latency percentile estimation |
| JSON and CSV | Exporting test results |

## How It Works

The application follows a modular workflow:

1. **CLI Parser:** Reads and validates the user's test configuration.
2. **Worker Pool:** Generates concurrent HTTP requests.
3. **Rate Limiter:** Regulates the rate at which requests are started.
4. **HTTP Client:** Sends requests and measures response latency.
5. **Metrics Collector:** Aggregates request counts, status codes, latency, and errors.
6. **Percentile Engine:** Calculates exact percentiles for validation and estimates percentiles using t-digest.
7. **Dashboard and Report Exporter:** Displays live results and saves the final report.

## Performance Metrics

- **Concurrency:** Maximum number of requests permitted to be in progress.
- **Request rate:** Target number of request starts per second.
- **Throughput:** Number of completed requests per second.
- **Latency:** Time taken to complete an HTTP request.
- **p50:** Median response latency.
- **p95:** Latency at or below which approximately 95% of measured requests complete.
- **p99:** Latency at or below which approximately 99% of measured requests complete.
- **Error rate:** Percentage of requests classified as errors under the application's defined counting rules.

The configured request rate and actual throughput are different. The achieved performance depends on factors such as server response time, concurrency, available resources, and network conditions.

## Requirements

- Go installed on the system.
- Visual Studio Code or another code editor (optional).
- A terminal or command prompt.
- Access to a local API or another API you are authorized to test.

Verify that Go is installed:

```bash
go version
```

## Setup and Installation

### 1. Open the Project

Open the APIStress project directory in Visual Studio Code.

Open the integrated terminal using **Terminal → New Terminal**.

### 2. Download Dependencies

Run the following command from the project root:

```bash
go mod download
```

If dependencies need to be resolved or `go.sum` needs to be generated, use:

```bash
go mod tidy
```

### 3. Build the Application

```bash
go build ./...
```

This checks whether the Go packages compile. If a standalone executable is required, use the appropriate build command for the project's CLI entry point.

### 4. Start the Local Demo Server

If the project includes a local demonstration server at `cmd/demo-server/main.go`, open a terminal and run:

```bash
go run ./cmd/demo-server
```

Keep this terminal open. Confirm the listening port and endpoint from the server's startup message or source code.

### 5. Run APIStress

Open a second terminal. Assuming the CLI entry point is `cmd/apistress/main.go`, run:

```bash
go run ./cmd/apistress \
  --url http://localhost:8080/api \
  --concurrency 50 \
  --duration 10s \
  --rate 100
```

This example requests a 10-second test with 50 concurrent workers and a target rate of 100 request starts per second.

**Important:** The example assumes that the project uses the specified entry point, port, endpoint, and flag names. Check the actual CLI implementation and demo-server configuration if your commands differ.

Use the help command to see the supported options:

```bash
go run ./cmd/apistress --help
```

## Running Tests

Run all Go tests:

```bash
go test ./...
```

Run tests with the Go race detector:

```bash
go test -race ./...
```

Run static analysis:

```bash
go vet ./...
```

A successful build or test run should be confirmed from the actual command output.

## Reports

The application is designed to export test results in two formats:

- **JSON:** Structured performance data for other programs and analysis tools.
- **CSV:** Tabular results that can be opened in spreadsheet applications.


