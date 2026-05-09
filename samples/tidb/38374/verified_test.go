// Race-trigger test for tidb-38374; see README.md for usage.

package domain

import (
	"sync"
	"testing"

	"tidb/sessionctx/variable"
)

func TestRace_tidb_38374(t *testing.T) {
	const ITERS = 5000
	noopA := func(a, b string) {}
	noopB := func(a, b string) {}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < ITERS; i++ {
			if i%2 == 0 { initSysVars(noopA) } else { initSysVars(noopB) }
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < ITERS; i++ {
			variable.SetGlobal("tidb_tso_client_batch_max_wait_time", "1.0")
		}
	}()
	wg.Wait()
}
