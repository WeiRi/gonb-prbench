package downloader

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX-version test. Downloader.mode is `uint32`.
// Accessed via atomic.LoadUint32 / StoreUint32 → no race.
func TestRace_go_ethereum_21201_mode(t *testing.T) {
	d := &Downloader{}
	var done atomic.Int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && done.Load() == 0; i++ {
			_ = atomic.LoadUint32(&d.mode)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && done.Load() == 0; i++ {
			atomic.StoreUint32(&d.mode, uint32(FullSync))
		}
		done.Store(1)
	}()
	wg.Wait()
}
