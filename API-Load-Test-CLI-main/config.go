package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	URL         string
	Concurrency int
	Duration    time.Duration
	Rate        int
	JSONOutput  string
	CSVOutput   string
}

func ParseConfig() (Config, error) {
	var cfg Config
	var duration string

	flag.StringVar(&cfg.URL, "url", "", "Target API URL")
	flag.IntVar(&cfg.Concurrency, "c", 10, "Number of concurrent workers")
	flag.StringVar(&duration, "d", "30s", "Test duration")
	flag.IntVar(&cfg.Rate, "r", 0, "Target request starts per second (0 = unlimited)")

	flag.StringVar(
		&cfg.JSONOutput,
		"json",
		"results/report.json",
		"Path for JSON report",
	)
	flag.StringVar(
		&cfg.CSVOutput,
		"csv",
		"results/report.csv",
		"Path for CSV report",
	)

	flag.Parse()

	if err := validateURL(cfg.URL); err != nil {
		return cfg, err
	}

	if cfg.Concurrency <= 0 {
		return cfg, errors.New("concurrency must be greater than 0")
	}

	if cfg.Rate < 0 {
		return cfg, errors.New("rate cannot be negative")
	}

	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		return cfg, fmt.Errorf("invalid duration: %w", err)
	}

	if parsedDuration <= 0 {
		return cfg, errors.New("duration must be greater than 0")
	}

	cfg.Duration = parsedDuration
	return cfg, nil
}

func validateURL(raw string) error {
	if raw == "" {
		return errors.New("URL is required; use -url http://localhost:8080/")
	}

	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if (parsed.Scheme != "http" && parsed.Scheme != "https") ||
		parsed.Host == "" {
		return errors.New("URL must be an absolute HTTP or HTTPS URL")
	}

	return nil
}
