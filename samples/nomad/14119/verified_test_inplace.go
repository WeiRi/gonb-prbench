// In-place race test for nomad-14119: package=client, uses upstream heartbeatStop.
// Bug: heartbeatstop.go:73 — watch() writes h.lastOk without lock,
// while setLastOk() writes under lock. Concurrent writes on bare field.
// PR fix: replace `h.lastOk = time.Now()` with `h.setLastOk(time.Now())`.
package client

import (
	"sync"
	"testing"
	"time"
)

func TestRace_14119_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 200

	for trial := 0; trial < ITERS; trial++ {
		h := &heartbeatStop{
			lock:        &sync.RWMutex{},
			shutdownCh:  make(chan struct{}),
			getRunner:   func(id string) (AllocRunner, error) { return nil, nil },
		}

		// Start watch() goroutine — writes lastOk without lock (heartbeatstop.go:73)
		go h.watch()

		var wg sync.WaitGroup

		// Concurrent setLastOk() calls — writes lastOk under lock (heartbeatstop.go:126-129)
		for i := 0; i < N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					h.setLastOk(time.Now())
				}
			}()
		}

		time.Sleep(5 * time.Millisecond)
		close(h.shutdownCh)
		wg.Wait()
	}
}
