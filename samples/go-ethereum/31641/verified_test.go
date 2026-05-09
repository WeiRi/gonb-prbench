package legacypool

import (
	"sync"
	"testing"
)

// TestRace_31641_legacypool_Add_vs_Clear: races LegacyPool.Add (reads pool.all)
// vs LegacyPool.Clear (writes pool.all). PR #31641 fixes by making Clear
// mutate the existing lookup in-place (with internal lock) rather than swap
// the pointer.
func TestRace_31641_legacypool_Add_vs_Clear(t *testing.T) {
	pool := NewLegacyPool()
	txs := []*Transaction{{hash: Hash{1}}, {hash: Hash{2}}, {hash: Hash{3}}}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			pool.Add(txs)
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
