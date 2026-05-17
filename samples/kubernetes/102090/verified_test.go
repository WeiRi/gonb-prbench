package testing

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Race: BUG Fake.AddWatchReactor / PrependWatchReactor append to
// c.WatchReactionChain WITHOUT holding c.Mutex. Concurrent goroutines
// race on slice append (header write).
// FIX adds c.Lock/Unlock around the append.
func TestRace_kubernetes_102090_fake_watchreactor(t *testing.T) {
	c := &Fake{}

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
			c.AddWatchReactor("v1", nil)
		}
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
			c.PrependWatchReactor("v2", nil)
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
