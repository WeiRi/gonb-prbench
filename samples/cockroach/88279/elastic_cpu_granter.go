package admission

import "sync"

// Reproduction of PR cockroachdb/cockroach#88279:
// "admission: squash data race in accessing token bucket"
// BUG: tb is at struct level; tryGet/tookWithoutPermission/setUtilizationLimit
// access tb without/with mu in inconsistent ways → race.

type TokenBucket struct {
	tokens int64
	rate   int64
}

func (b *TokenBucket) TryToFulfill(n int64) bool {
	if b.tokens >= n {
		b.tokens -= n
		return true
	}
	return false
}

func (b *TokenBucket) Adjust(delta int64) { b.tokens += delta }
func (b *TokenBucket) UpdateConfig(rate, burst int64) { b.rate = rate; b.tokens = burst }

type elasticCPUGranter struct {
	mu struct {
		sync.Mutex
		utilizationLimit float64
	}
	tb *TokenBucket // BUG: not under mu in pre-fix code
}

// tryGet (BUG): no lock around tb access.
func (e *elasticCPUGranter) tryGet(count int64) bool {
	return e.tb.TryToFulfill(count) // line 36
}

// tookWithoutPermission (BUG): no lock.
func (e *elasticCPUGranter) tookWithoutPermission(count int64) {
	e.tb.Adjust(-count) // line 41
}

// setUtilizationLimit (BUG): mutates e.tb under mu but tryGet/took read without mu.
func (e *elasticCPUGranter) setUtilizationLimit(rate float64) {
	e.mu.Lock()
	e.mu.utilizationLimit = rate
	e.mu.Unlock()
	// BUG: tb.UpdateConfig outside the lock and tb is a pointer that can be
	// concurrently re-read by tryGet.
	e.tb.UpdateConfig(int64(rate*1e9), int64(rate*1e9)) // line 51
}

func newElasticCPUGranter() *elasticCPUGranter {
	return &elasticCPUGranter{tb: &TokenBucket{tokens: 1000, rate: 1e9}}
}

