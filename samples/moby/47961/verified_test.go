package client

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG-version test. Field `negotiated` is `bool` (plain).
// Concurrent unsynchronized read/write fires race detector.
func TestRace_moby_47961_negotiated(t *testing.T) {
	cli := &Client{negotiateVersion: true}
	var done atomic.Int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && done.Load() == 0; i++ {
			_ = cli.negotiated
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && done.Load() == 0; i++ {
			cli.negotiated = true
		}
		done.Store(1)
	}()
	wg.Wait()
}
