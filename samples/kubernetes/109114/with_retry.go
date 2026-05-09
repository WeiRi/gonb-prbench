package rest

import (
	"context"
	"net/url"
	"time"
)

// Stripped reproduction of staging/src/k8s.io/client-go/rest/with_retry.go pre-PR #109114.
// BUG: Request.retry is a shared mutable *withRetry instance; callers race on its fields.

type RESTClient struct {
	base *url.URL
}

type Request struct {
	c     *RESTClient
	retry *withRetry
}

type withRetry struct {
	maxRetries int
	attempts   int
	retryAfter time.Duration
}

// Before — BUG: writes retry fields without lock.
func (r *withRetry) Before(ctx context.Context, req *Request) {
	r.attempts++                          // line 28 — racing write
	r.retryAfter = 0                      // line 29 — racing write
}

// After — BUG: reads retry fields without lock; also writes attempts.
func (r *withRetry) After(ctx context.Context, req *Request, _ interface{}, _ error) {
	_ = r.attempts                        // line 34 — racing read
	if r.attempts >= r.maxRetries {       // line 35 — racing read
		return
	}
	r.retryAfter = time.Millisecond       // line 38 — racing write
}
