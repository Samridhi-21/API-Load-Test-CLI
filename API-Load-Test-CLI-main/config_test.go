package main

import (
	"testing"
	"time"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid HTTP", "http://localhost:8080/", false},
		{"valid HTTPS", "https://example.com/api", false},
		{"empty URL", "", true},
		{"missing scheme", "localhost:8080", true},
		{"unsupported scheme", "ftp://example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateURL(%q) error = %v, wantErr %v",
					tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestDurationParsing(t *testing.T) {
	got, err := time.ParseDuration("30s")
	if err != nil {
		t.Fatal(err)
	}

	if got != 30*time.Second {
		t.Fatalf("got %v, want 30s", got)
	}
}
