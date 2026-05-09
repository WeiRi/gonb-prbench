// Race-trigger test for tidb-40277; see README.md for usage.

package expression_verified_tidb40277

import (
	"sync"
	"testing"
)

// StatementContext models the tidb StmtCtx field that the bug shared.
type StatementContext struct {
	OverflowAsWarning bool
	IgnoreTruncate    bool
	TruncateAsWarning bool
	InInsertStmt      bool
}

type SessionVars struct{ StmtCtx *StatementContext }

type ExprCtx struct{ vars *SessionVars }

func (c *ExprCtx) GetSessionVars() *SessionVars { return c.vars }

type castJSONAsArrayFunctionSig struct{ ctx *ExprCtx }

func newSig(sc *StatementContext) *castJSONAsArrayFunctionSig {
	return &castJSONAsArrayFunctionSig{ctx: &ExprCtx{vars: &SessionVars{StmtCtx: sc}}}
}

// EvalJSON_BUG is the pre-PR40277 form (expression/builtin_cast.go:478..489):
// it MUTATES three fields on the shared sc and restores via defer.
func (b *castJSONAsArrayFunctionSig) EvalJSON_BUG(vals []int) (sum int) {
	sc := b.ctx.GetSessionVars().StmtCtx
	originalOverflowAsWarning := sc.OverflowAsWarning
	originIgnoreTruncate := sc.IgnoreTruncate
	originTruncateAsWarning := sc.TruncateAsWarning
	sc.OverflowAsWarning = false // race line ~478
	sc.IgnoreTruncate = false
	sc.TruncateAsWarning = false
	defer func() {
		sc.OverflowAsWarning = originalOverflowAsWarning // race line ~481
		sc.IgnoreTruncate = originIgnoreTruncate
		sc.TruncateAsWarning = originTruncateAsWarning
	}()
	for _, v := range vals {
		sum += v
		if sc.OverflowAsWarning { // race line ~489 (read)
			sum++
		}
	}
	return
}

// TestRace_tidb40277 reproduces the data race by running two goroutines that
// concurrently evaluate over the SAME StatementContext (same session).
func TestRace_tidb40277(t *testing.T) {
	sc := &StatementContext{}
	sigA := newSig(sc)
	sigB := newSig(sc)

	vals := []int{1, 2, 3, 4, 5}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			_ = sigA.EvalJSON_BUG(vals)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			_ = sigB.EvalJSON_BUG(vals)
		}
	}()
	wg.Wait()
}
