// Minimal extraction of pre-fix go-ethereum core/txpool/legacypool/legacypool.go
// for PR #31641. Pre-fix race: LegacyPool.Clear assigns `pool.all = newLookup()`
// while Add reads `pool.all.Get(...)` concurrently — pointer-replacement race.
// Production-code path: core/txpool/legacypool/legacypool.go
package legacypool

import "sync"

type Hash [32]byte

type Transaction struct {
	hash Hash
}

func (tx *Transaction) Hash() Hash { return tx.hash }

// lookup mirrors the upstream type with internal lock + map.
type lookup struct {
	lock  sync.RWMutex
	slots int
	txs   map[Hash]*Transaction
}

func newLookup() *lookup {
	return &lookup{txs: make(map[Hash]*Transaction)}
}

func (t *lookup) Get(h Hash) *Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.txs[h]
}

func (t *lookup) Add(tx *Transaction) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.txs[tx.hash] = tx
	t.slots++
}

// LegacyPool — pre-fix subset.
type LegacyPool struct {
	mu  sync.RWMutex
	all *lookup
}

func NewLegacyPool() *LegacyPool {
	return &LegacyPool{all: newLookup()}
}

// Add — pre-fix simplified path: reads pool.all.Get(tx.Hash()) WITHOUT holding pool.mu.
// Upstream path: core/txpool/legacypool/legacypool.go (pre-fix line ~947).
func (pool *LegacyPool) Add(txs []*Transaction) {
	for _, tx := range txs {
		if pool.all.Get(tx.Hash()) != nil {
			continue
		}
		pool.all.Add(tx)
	}
}

// Clear — pre-fix version: REASSIGNS pool.all = newLookup() unprotected by pool.mu.
// PR fix replaces this with `pool.all.Clear()` (no pointer swap).
// Upstream path: core/txpool/legacypool/legacypool.go (pre-fix line ~1929).
func (pool *LegacyPool) Clear() {
	pool.all = newLookup()
}
