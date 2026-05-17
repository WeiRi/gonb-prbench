package stmtctx

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: StmtCtx.IsDDLJobInQueue is a plain bool. DDLExec.Next defers
// `IsDDLJobInQueue = false` while InfoSyncer.ReportMinStartTS concurrently
// reads `IsDDLJobInQueue` -> race on plain bool.
func TestRace_tidb_62900_stmtctx_isddljobinqueue(t *testing.T) {
	sc := &StatementContext{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 200000 && atomic.LoadInt32(&done) == 0; j++ {
			sc.IsDDLJobInQueue = true
			sc.IsDDLJobInQueue = false
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 200000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = sc.IsDDLJobInQueue
		}
	}()
	wg.Wait()
}
