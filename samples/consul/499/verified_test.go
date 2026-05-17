package consul

import (
	"io"
	"sync"
	"testing"
	"time"
)

// TestRace_Consul499_ReapShutdown triggers the race between ConnPool.reap()
// reading p.shutdown without lock (pool.go:361) and Shutdown() writing
// p.shutdown under p.Lock().
//
// Fix: changed reap() from "for !p.shutdown" to "for { select { case
// <-p.shutdownCh: return; case <-time.After(time.Second): } }"
//
// NOTE: Requires proper consul dependency resolution (go mod tidy with
// compatible dependency versions from the consul-499 era).
func TestRace_Consul499_ReapShutdown(t *testing.T) {
	const numGoroutines = 50
	const iterations = 200

	for i := 0; i < iterations; i++ {
		// NewPool with maxTime > 0 launches reap() goroutine automatically
		p := NewPool(io.Discard, time.Hour, 5, nil)

		var wg sync.WaitGroup

		// Concurrent Shutdown calls write p.shutdown=true under p.Lock(),
		// racing with reap's unlocked read of p.shutdown
		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = p.Shutdown()
			}()
		}

		wg.Wait()
	}
}
