package legacypool

import (
	"sync"
	"testing"
)

// TestRace_31758_legacypool_pricedList_Clear_vs_read: races readers of
// pool.priced vs LegacyPool.Clear which reassigns pool.priced.
func TestRace_31758_legacypool_pricedList_Clear_vs_read(t *testing.T) {
	pool := NewLegacyPool()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			_ = pool.Probe()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			pool.Clear()
		}
	}()
	wg.Wait()
}
