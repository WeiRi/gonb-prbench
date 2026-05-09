package filters

import (
	"net/http"
	"sync"
	"testing"
)

// TestRace_130741_ProductionPath reproduces the data race from PR #130741.
//
// The bug: package-level int32 vars (atomicMutatingExecuting, atomicReadOnlyExecuting,
// atomicMutatingWaiting, atomicReadOnlyWaiting) are written via atomic.AddInt32 in
// priority-and-fairness.go:146-148 during request processing. But test code reads them
// directly without atomic.LoadInt32. This causes a data race when the test goroutine
// reads these vars while a handler goroutine writes them via atomic.AddInt32.
//
// The fix changes int32 -> atomic.Int32 and uses .Load() for all reads.
func TestRace_130741_ProductionPath(t *testing.T) {
	// Use decisionNoQueuingExecute: the handler WILL execute through the APF path,
	// which calls noteExecutingDelta -> atomic.AddInt32 in priority-and-fairness.go:146-148
	server := newApfServerWithSingleRequest(t, decisionNoQueuingExecute)
	defer server.Close()

	var wg sync.WaitGroup
	const N = 300

	// Fire concurrent HTTP requests through the real APF handler.
	// Each request triggers atomic.AddInt32 writes in priority-and-fairness.go:146-148.
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < N; i++ {
				resp, err := http.Get(server.URL + "/api/v1/namespaces/default/pods")
				if err == nil {
					resp.Body.Close()
				}
			}
		}()
	}

	// Concurrently perform non-atomic reads of the same package-level variables.
	// These DIRECT reads RACE with atomic.AddInt32 writes from handler goroutines.
	for g := 0; g < 5; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < N*2; i++ {
				_ = atomicMutatingExecuting   // non-atomic read - RACE with atomic writes
				_ = atomicReadOnlyExecuting   // non-atomic read - RACE
				_ = atomicMutatingWaiting     // non-atomic read - RACE
				_ = atomicReadOnlyWaiting     // non-atomic read - RACE
			}
		}()
	}

	wg.Wait()
}
