package main

import (
	"math"
	"sort"
	"time"
)

// Percentile calculates the nearest-rank percentile.
// p must be greater than 0 and at most 100.
// It returns false when the input is empty or p is invalid.
func Percentile(
	latencies []time.Duration,
	p float64,
) (time.Duration, bool) {
	if len(latencies) == 0 || p <= 0 || p > 100 {
		return 0, false
	}

	// Sort a copy so the caller's slice remains unchanged.
	sorted := append([]time.Duration(nil), latencies...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	rank := int(math.Ceil(p / 100 * float64(len(sorted))))
	return sorted[rank-1], true
}
