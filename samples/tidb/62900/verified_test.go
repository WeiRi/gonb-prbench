// Race-trigger test for tidb-62900; see README.md for usage.

package ddl_verified_tidb62900

import (
	"sync"
	"testing"
)

// StatementContext mirrors the tidb StmtCtx field (bool form, pre-PR62900).
type StatementContext struct {
	IsDDLJobInQueue bool
}

type SessionVars struct{ StmtCtx *StatementContext }

type Executor struct{ vars *SessionVars }

// DoDDLJobWrapper_BUG is pkg/ddl/executor.go:6771 pre-fix — direct assign.
func (e *Executor) DoDDLJobWrapper_BUG() {
	e.vars.StmtCtx.IsDDLJobInQueue = true
}

// ReportMinStartTS_BUG is pkg/domain/infosync/info.go:630 pre-fix — direct read.
func ReportMinStartTS_BUG(sc *StatementContext) bool {
	return sc.IsDDLJobInQueue
}

func TestRace_tidb62900(t *testing.T) {
	vars := &SessionVars{StmtCtx: &StatementContext{}}
	exec := &Executor{vars: vars}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			exec.DoDDLJobWrapper_BUG()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = ReportMinStartTS_BUG(vars.StmtCtx)
		}
	}()
	wg.Wait()
}
