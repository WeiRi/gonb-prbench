package agent

import (
	"sync"
	"testing"
)

// TestShutdownRace exercises concurrent access to the shutdown fields in Agent.
// PR 3162: fix 'consul leave' shutdown race where HTTP server shutdown inside
// Shutdown() blocked delivery of the HTTP response.
// Original diff file: agent/agent.go
// Original frame hits: agent/agent.go:1097 (Shutdown)
func TestShutdownRace(t *testing.T) {
	// Minimal agent-like struct replicating the shutdown race pattern
	type miniAgent struct {
		shutdownLock sync.Mutex
		shutdown     bool
		shutdownCh   chan struct{}
	}

	numGoroutines := 60
	iterations := 200

	for i := 0; i < iterations; i++ {
		a := &miniAgent{
			shutdownCh: make(chan struct{}),
		}

		var wg sync.WaitGroup

		// Multiple goroutines calling Shutdown concurrently
		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				a.shutdownLock.Lock()
				if a.shutdown {
					a.shutdownLock.Unlock()
					return
				}
				a.shutdown = true
				a.shutdownLock.Unlock()
				close(a.shutdownCh)
			}()
		}

		// Goroutine reading shutdown without lock (simulating HTTP handler check)
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Read shutdown WITHOUT lock - RACY
			_ = a.shutdown
			select {
			case <-a.shutdownCh:
			default:
			}
		}()

		wg.Wait()
	}
}
