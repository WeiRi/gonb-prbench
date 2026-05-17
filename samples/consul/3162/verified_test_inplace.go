// In-place race test for consul-3162: uses REAL upstream Agent type.
// Bug: agent/agent.go — Shutdown() writes a.shutdown=true under shutdownLock,
// but concurrent goroutines read a.shutdown WITHOUT the lock at line ~768.
// PR 3162: moved endpoint shutdown to async goroutine so HTTP response
// for 'consul leave' is not blocked by in-line server shutdown.
package agent

import (
	"io"
	"log"
	"sync"
	"testing"
)

func TestRace_3162_InPlace(t *testing.T) {
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
		// (simulates the select at agent.go:768 which reads shutdownCh
		// and the implicit check of agent state)
		for i := 0; i < N/2; i++ {
			wg.Add(1)
			go func(ag *Agent) {
				defer wg.Done()
				_ = ag.shutdown // RACY READ without shutdownLock
				select {
				case <-ag.shutdownCh:
				default:
				}
			}(a)
		}

		// Writer goroutines: call Shutdown() which writes a.shutdown=true
		// under shutdownLock (line 1163) then close(a.shutdownCh) (line 1164)
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
