package expression

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/pingcap/tidb/sessionctx/stmtctx"
)

// BUG: castJSONAsArrayFunctionSig.evalJSON mutates shared
// b.ctx.GetSessionVars().StmtCtx fields (OverflowAsWarning etc); restored
// via defer. Concurrent evalJSON calls share same StmtCtx → race on these fields.
func TestRace_tidb_40277_stmt_ctx_fields(t *testing.T) {
	sc := &stmtctx.StatementContext{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			orig := sc.OverflowAsWarning
			sc.OverflowAsWarning = false
			sc.OverflowAsWarning = orig
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = sc.OverflowAsWarning
		}
	}()
	wg.Wait()
}
