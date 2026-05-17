package downloader

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG-version test. Downloader.mode is `SyncMode` (uint, plain).
// Concurrent unsynchronized read/write fires race detector.
func TestRace_go_ethereum_21201_mode(t *testing.T) {
	d := &Downloader{}
	var done atomic.Int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && done.Load() == 0; i++ {
			_ = d.mode
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && done.Load() == 0; i++ {
			d.mode = FullSync
		}
		done.Store(1)
	}()
	wg.Wait()
}
