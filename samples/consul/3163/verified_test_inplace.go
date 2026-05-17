// In-place race test for consul-3163: uses REAL upstream Agent type.
// Bug: agent/agent.go — split Shutdown into ShutdownAgent/ShutdownEndpoints
// to properly sequence the shutdown with mutex.
// The race: Shutdown() writes a.shutdown=true under shutdownLock at ~line 1163,
// while concurrent readers access a.shutdown without the lock.
package agent

import (
	"io"
	"log"
	"sync"
	"testing"
)

func TestRace_3163_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 100

	for iter := 0; iter < ITERS; iter++ {
		a := &Agent{
			config:     &Config{},
			logger:     log.New(io.Discard, "", 0),
			shutdownCh: make(chan struct{}),
		}

		var wg sync.WaitGroup

		// Reader goroutines: read a.shutdown WITHOUT holding shutdownLock
		for i := 0; i < N/2; i++ {
			wg.Add(1)
			go func(ag *Agent) {
				defer wg.Done()
				_ = ag.shutdown // RACY READ without lock
				select {
				case <-ag.shutdownCh:
				default:
				}
			}(a)
		}

		// Writer goroutines: call Shutdown() which WRITES a.shutdown under lock
		for i := 0; i < N/2; i++ {
			wg.Add(1)
			go func(ag *Agent) {
				defer wg.Done()
				defer func() { recover() }()
				_ = ag.Shutdown()
			}(a)
		}

		wg.Wait()
	}
}
