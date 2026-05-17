package client

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX-version test. Field `negotiated` is `atomic.Bool`.
// Calls Load/Store which are atomic → no race.
func TestRace_moby_47961_negotiated(t *testing.T) {
	cli := &Client{negotiateVersion: true}
	var done atomic.Int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && done.Load() == 0; i++ {
			_ = cli.negotiated.Load()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && done.Load() == 0; i++ {
			cli.negotiated.Store(true)
		}
		done.Store(1)
	}()
	wg.Wait()
}
