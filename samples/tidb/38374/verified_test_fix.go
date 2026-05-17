// Race test for tidb-38374 FIX state (atomic.Pointer wrapper)
// Verifies atomic.Pointer API is race-free.
package variable

import (
	"sync"
	"testing"
)

func TestRace_38374_FIX(t *testing.T) {
	var wg sync.WaitGroup
	const N = 40
	const ITERS = 500
	for i := 0; i < N; i++ {
		wg.Add(2)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				f := func(a, b string) {}
				SetPDClientDynamicOption.Store(&f)
			}
		}(i)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				if p := SetPDClientDynamicOption.Load(); p != nil {
					(*p)("x", "y")
				}
			}
		}()
	}
	wg.Wait()
}
