// PR: https://github.com/moby/moby/pull/26695
// Fix: Add sync.Mutex to pauseMonitor to prevent data race on waiters map.
package libcontainerd

import (
	"sync"
	"testing"
)

func TestRace26695PauseMonitor(t *testing.T) {
	pm := &pauseMonitor{}

	var wg sync.WaitGroup
	const iterations = 1000
	wg.Add(iterations * 2)

	// Concurrent handle (reads/clears waiters) and append (adds waiters).
	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()
			ch := make(chan struct{})
			pm.append("test", ch)
			pm.handle("test")
		}()
		go func() {
			defer wg.Done()
			ch := make(chan struct{})
			pm.append("other", ch)
			pm.handle("other")
		}()
	}

	wg.Wait()
}
