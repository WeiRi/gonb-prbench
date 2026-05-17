package expression

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/pingcap/tidb/sessionctx/stmtctx"
)

// FIX: each call uses fakeSctx (a dedicated StatementContext) instead of
// mutating the shared b.ctx.StmtCtx. No race on shared fields.
func TestRace_tidb_40277_stmt_ctx_fields(t *testing.T) {
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			fakeSctx := &stmtctx.StatementContext{InInsertStmt: true}
			fakeSctx.OverflowAsWarning = false
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			fakeSctx := &stmtctx.StatementContext{InInsertStmt: true}
			_ = fakeSctx.OverflowAsWarning
		}
	}()
	wg.Wait()
}
