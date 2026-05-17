package legacypool

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Race: BUG pool.Clear() reassigns pool.priced = newPricedList(pool.all)
// while concurrent readers of pool.priced FIELD race on the pointer write.
// FIX: pool.Clear() calls pool.priced.Reheap() preserving the pointer.
func TestRace_go_ethereum_31758_pool_priced(t *testing.T) {
	pool, _ := setupPool()
	defer pool.Close()

	var done atomic.Bool
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 200 && !done.Load(); i++ {
			pool.Clear()
		}
		done.Store(true)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 500000 && !done.Load(); i++ {
			_ = pool.priced
		}
	}()
	wg.Wait()
}
