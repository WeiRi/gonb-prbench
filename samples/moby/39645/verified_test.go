package container

import (
	"sync"
	"testing"
)

func TestRace_39645(t *testing.T) {
	const N = 50
	const ITERS = 200

	h := &Health{}

	var wg sync.WaitGroup
	wg.Add(N * 2)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				h.mu.Lock()
				h.Health.Status = "healthy"
				h.mu.Unlock()
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = h.Health.Status
			}
		}()
	}
	wg.Wait()
}
