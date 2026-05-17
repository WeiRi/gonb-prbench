package legacypool

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// Race: BUG pool.Clear() reassigns pool.all = newLookup() under pool.mu, but
// concurrent readers of pool.all (no pool.mu lock) race on the pointer.
// FIX: pool.Clear() calls pool.all.Clear() preserving the pointer; lookup's
// own RWMutex protects map writes.
func TestRace_go_ethereum_31641_pool_all(t *testing.T) {
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
		h := common.Hash{}
		for i := 0; i < 200000 && !done.Load(); i++ {
			_ = pool.all.Get(h)
		}
	}()
	wg.Wait()
}
