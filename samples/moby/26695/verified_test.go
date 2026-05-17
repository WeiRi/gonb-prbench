// Race test for moby-26695 (libcontainerd.pauseMonitor)
// Targets the race in pausemonitor_linux.go:
//   - handle() reads & clears m.waiters map
//   - append() adds to m.waiters map
// BUG: no lock. FIX: sync.Mutex added.
package libcontainerd

import (
	"sync"
	"testing"
)

func TestRace_26695_PauseMonitor(t *testing.T) {
	pm := &pauseMonitor{}

	var wg sync.WaitGroup
	const iterations = 1000

	for i := 0; i < iterations; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			ch := make(chan struct{})
			pm.append("test", ch)
		}()
		go func() {
			defer wg.Done()
			pm.handle("test")
		}()
	}
	wg.Wait()
}
