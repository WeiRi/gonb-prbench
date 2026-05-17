package agent

import (
	"sync"
	"testing"
)

// BUG: runManager's `go func(){...<-ready... }()` closes over outer `ready`
// variable. After loop iteration, `ready = nil` mutates the captured var,
// racing with the goroutine reading <-ready.
//
// Synthetic test that mimics the production race pattern.
func TestRace_swarmkit_1046_closure_capture(t *testing.T) {
	const N = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		ready := make(chan struct{})
		go func() {
			defer wg.Done()
			select {
			case <-ready:
			default:
			}
		}()
		go func() {
			defer wg.Done()
			ready = nil
		}()
	}
	wg.Wait()
}
