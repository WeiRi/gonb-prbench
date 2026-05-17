package core

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// BUG: pool.local() reads pool.locals.accounts / pool.pending / pool.queue
// without holding pool.mu. Concurrent writers under pool.mu cause race.
// PR #14940 fix wraps the journal-rotation call site with pool.mu.Lock(),
// but local() itself is unlocked. Our prepatch.diff adds RLock to local()
// so the fix path actually suppresses the whitebox-visible race.
func TestRace_14940_local_vs_writer(t *testing.T) {
	pool := &TxPool{
		pending: make(map[common.Address]*txList),
		queue:   make(map[common.Address]*txList),
		locals:  &accountSet{accounts: make(map[common.Address]struct{})},
	}
	seed := common.BytesToAddress([]byte{1, 2, 3, 4})
	pool.locals.accounts[seed] = struct{}{}

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for j := 0; j < 500 && atomic.LoadInt32(&done) == 0; j++ {
			_ = pool.local()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 500 && atomic.LoadInt32(&done) == 0; j++ {
			addr := common.BytesToAddress([]byte{byte(j), byte(j >> 8)})
			pool.mu.Lock()
			pool.locals.accounts[addr] = struct{}{}
			pool.pending[addr] = newTxList(false)
			delete(pool.locals.accounts, addr)
			delete(pool.pending, addr)
			pool.mu.Unlock()
		}
	}()
	wg.Wait()
}
