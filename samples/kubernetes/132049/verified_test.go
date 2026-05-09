package handlers

import (
	"context"
	"sync"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/finisher"
)

// TestRace_132049_TimeoutWasCreated triggers the data race where
// patchResource reads wasCreated after finisher.FinishRequest returns
// a timeout error, while the goroutine executing the request is still
// running and may write to wasCreated (captured in the closure).
//
// The fix (in patch.go) adds a check for errors.IsTimeout before
// reading wasCreated, returning false instead to avoid the race.
func TestRace_132049_TimeoutWasCreated(t *testing.T) {
	numGoroutines := 50
	iterations := 200
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				var wasCreated bool

				// Simulate the requestFunc closure in patchResource
				// that sets wasCreated = created (patch.go:705)
				requestFunc := func() (runtime.Object, error) {
					wasCreated = true
					return nil, nil
				}

				// Simulate timeout via cancelled context
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				// finisher.FinishRequest returns timeout error while
				// the goroutine calling requestFunc is still running
				result, err := finisher.FinishRequest(ctx, requestFunc)

				// In BUG state: reading wasCreated after timeout
				// races with goroutine still writing to it
				_ = result
				_ = wasCreated // DATA RACE with requestFunc goroutine
				_ = err
			}
		}()
	}

	wg.Wait()
}
