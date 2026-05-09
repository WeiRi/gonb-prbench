// Minimal extraction of pre-fix go-ethereum core/txpool/legacypool/legacypool.go
// for PR #31758. Pre-fix race: LegacyPool.Clear assigns
// `pool.priced = newPricedList(pool.all)` while concurrent readers (e.g.
// price-discounting paths) read pool.priced without lock — pointer-replacement race.
// PR fix replaces the assignment with `pool.priced.Reheap()`.
// Production-code path: core/txpool/legacypool/legacypool.go
package legacypool

import "sync"

type lookup struct {
	lock sync.RWMutex
}

func newLookup() *lookup { return &lookup{} }

// pricedList stub mirrors the upstream behaviour relevant to the race:
// holds a reference to all txs and exposes a Reheap method.
type pricedList struct {
	all  *lookup
	heap []int
}

func newPricedList(all *lookup) *pricedList {
	return &pricedList{all: all, heap: make([]int, 0, 8)}
}

func (p *pricedList) Reheap() {
	p.heap = p.heap[:0]
}

// Used returns a value derived from p.heap; serves as a read-side path of
// pool.priced from another goroutine.
func (p *pricedList) Used() int {
	return len(p.heap)
}

type LegacyPool struct {
	all    *lookup
	priced *pricedList
}

func NewLegacyPool() *LegacyPool {
	all := newLookup()
	return &LegacyPool{all: all, priced: newPricedList(all)}
}

// Probe simulates a hot path that reads pool.priced (e.g., during fee
// estimation or eviction). Pre-fix this read can race with Clear's reassignment.
// Upstream call sites read pool.priced from non-Clear goroutines.
func (pool *LegacyPool) Probe() int {
	return pool.priced.Used()
}

// Clear — pre-fix version: REASSIGNS pool.priced = newPricedList(pool.all).
// PR fix: pool.priced.Reheap() (no pointer swap).
// Upstream path: core/txpool/legacypool/legacypool.go (pre-fix line ~1937).
func (pool *LegacyPool) Clear() {
	pool.priced = newPricedList(pool.all)
}
