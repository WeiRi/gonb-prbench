package agent

import (
	"sync"
	"testing"
)

// FIX: goroutine captures `ready` as a parameter (closure-by-value).
// Outer mutation `ready = nil` no longer affects goroutine.
func TestRace_swarmkit_1046_closure_capture(t *testing.T) {
	const N = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		ready := make(chan struct{})
		go func(ready chan struct{}) {
			defer wg.Done()
			select {
			case <-ready:
			default:
			}
		}(ready)
		go func() {
			defer wg.Done()
			ready = nil
		}()
	}
	wg.Wait()
}
