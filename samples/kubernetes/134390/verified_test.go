package metrics

import (
	"sync"
	"testing"
)

func TestRace_134390(t *testing.T) {
	c := NewCounter(&CounterOpts{
		Namespace: "test_ns",
		Subsystem: "test_ss",
		Name:      "race_hidden",
		Help:      "test",
	})

	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				// Bug: ClearState writes isHidden/isDeprecated under createLock,
				// but IsHidden/IsDeprecated read them WITHOUT any lock.
				c.ClearState()
				_ = c.IsHidden()
				_ = c.IsDeprecated()
			}
		}()
	}

	wg.Wait()
}
