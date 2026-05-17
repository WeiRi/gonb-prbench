package stmtctx

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX: StmtCtx.IsDDLJobInQueue is atomic.Bool with Store/Load.
func TestRace_tidb_62900_stmtctx_isddljobinqueue(t *testing.T) {
	sc := &StatementContext{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 200000 && atomic.LoadInt32(&done) == 0; j++ {
			sc.IsDDLJobInQueue.Store(true)
			sc.IsDDLJobInQueue.Store(false)
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 200000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = sc.IsDDLJobInQueue.Load()
		}
	}()
	wg.Wait()
}
