package disk

import (
	"sync"
	"testing"
)

func TestRace_123215(t *testing.T) {
	const numGoroutines = 60
	const numIterations = 300

	for iter := 0; iter < numIterations; iter++ {
		m := NewMonitor()
		var wg sync.WaitGroup

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = m.IncrementalStats()
			}()
		}

		wg.Wait()
	}
}
