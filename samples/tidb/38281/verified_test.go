// Race test for tidb-38281 — collationInfo.coerInit plain bool race
// BUG: SetCoercibility writes c.coerInit = true without sync; HasCoercibility reads
// FIX: coerInit is atomic.Bool
package expression

import (
	"sync"
	"testing"
)

func TestRace_38281_CoerInit(t *testing.T) {
	c := &collationInfo{}
	const N = 50
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < N*20; i++ {
			c.SetCoercibility(CoercibilityExplicit)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N*20; i++ {
			_ = c.HasCoercibility()
		}
	}()
	wg.Wait()
}
