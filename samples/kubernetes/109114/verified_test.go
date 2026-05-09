package rest

import (
	"context"
	"net/url"
	"sync"
	"testing"
)

func TestRace_109114(t *testing.T) {
	r := &Request{
		c: &RESTClient{
			base: &url.URL{Scheme: "http", Host: "localhost"},
		},
		retry: &withRetry{maxRetries: 10},
	}

	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				// Bug: r.retry is a shared mutable WithRetry instance.
				// Before/After modify internal state (attempts, retryAfter).
				// Multiple goroutines executing requests on the same Request
				// race on the shared retry fields.
				ctx := context.Background()
				r.retry.After(ctx, r, nil, nil)
				r.retry.Before(ctx, r)
			}
		}()
	}

	wg.Wait()
}
