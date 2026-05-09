// PR #14940 - core/tx_pool.go - txpool journal/test races. Pre-fix: journal
// rotation case in loop() calls pool.journal.rotate(pool.local()) WITHOUT
// holding pool.mu, while another path (resetState) mutates pool.pending/local
// under pool.mu. PR fix: wrap journal rotation with pool.mu.Lock/Unlock.
// Production-code path: core/tx_pool.go (pre-fix line ~302).
package txpool

import "sync"

type TxPool struct {
	mu      sync.Mutex
	pending map[int]int
	locals  map[int]bool
}

func NewTxPool() *TxPool {
	return &TxPool{pending: make(map[int]int), locals: make(map[int]bool)}
}

// Pre-fix loop's journal-rotation branch calls Local() without lock.
// Upstream: core/tx_pool.go (pre-fix line ~302).
func (pool *TxPool) RotateJournal() []int {
	// PRE-FIX: NO pool.mu.Lock() here.
	out := make([]int, 0, len(pool.locals))
	for k := range pool.locals {
		if pool.pending[k] > 0 {
			out = append(out, k)
		}
	}
	return out
}

// reset mutates pool state under pool.mu.
func (pool *TxPool) Reset(addr int) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.pending[addr] = 1
	pool.locals[addr] = true
}
