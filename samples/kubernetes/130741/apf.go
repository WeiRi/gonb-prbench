// Production stub for kubernetes-130741.
// Pre-PR: package-level int32 vars are mutated by atomic.AddInt32 in handler
// path, but tests read them directly (non-atomic) — RACE.
package filters

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

// Pre-fix: plain int32 (not atomic.Int32) — RACE with atomic ops.
var (
	atomicMutatingExecuting int32
	atomicReadOnlyExecuting int32
	atomicMutatingWaiting   int32
	atomicReadOnlyWaiting   int32
)

type apfDecision int

const (
	decisionNoQueuingExecute apfDecision = iota
	decisionExecute
)

// noteExecutingDelta uses atomic.AddInt32 — racing with non-atomic reads.
func noteExecutingDelta(_ apfDecision) {
	atomic.AddInt32(&atomicMutatingExecuting, 1)
	atomic.AddInt32(&atomicReadOnlyExecuting, 1)
	atomic.AddInt32(&atomicMutatingWaiting, 1)
	atomic.AddInt32(&atomicReadOnlyWaiting, 1)
}

func newApfServerWithSingleRequest(_ *testing.T, dec apfDecision) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noteExecutingDelta(dec)
		w.WriteHeader(200)
	}))
}
