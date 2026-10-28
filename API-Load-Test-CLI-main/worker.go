package main

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

func RunWorkers(
	ctx context.Context,
	cfg Config,
	executor *HTTPExecutor,
) <-chan RequestResult {
	results := make(chan RequestResult, cfg.Concurrency*2)

	var wg sync.WaitGroup

	// Create ONE shared limiter for all workers.
	var limiter *rate.Limiter

	if cfg.Rate > 0 {
		// One token initially, then one token per interval.
		limiter = rate.NewLimiter(rate.Limit(cfg.Rate), 1)
	}

	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				if ctx.Err() != nil {
					return
				}

				// Wait until a token is available.
				if limiter != nil {
					if err := limiter.Wait(ctx); err != nil {
						return
					}
				}

				// Don't start another request after cancellation.
				if ctx.Err() != nil {
					return
				}

				result := executor.Do(ctx, cfg.URL)
				results <- result

				if ctx.Err() != nil {
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
