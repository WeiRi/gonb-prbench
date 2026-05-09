package miner

import (
	"sync"
	"testing"
)

// TestRace_17173_state_concurrent_map: Wait writes work.state.objects map
// without currentMu while Pending iterates the same map under currentMu — 
// concurrent map read/write fires under -race.
func TestRace_17173_state_concurrent_map(t *testing.T) {
	w := newWorker()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			w.Wait()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			_ = w.Pending()
		}
	}()
	wg.Wait()
}
