// Race test for tidb-38374 BUG state (plain func var)
// SetPDClientDynamicOption is declared as `func(string, string) = nil`
// Concurrent write (assign new func) + read (call/check) races.
package variable

import (
	"sync"
	"testing"
)

func TestRace_38374_BUG(t *testing.T) {
	var wg sync.WaitGroup
	const N = 40
	const ITERS = 500
	for i := 0; i < N; i++ {
		wg.Add(2)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				SetPDClientDynamicOption = func(a, b string) {}
			}
		}(i)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				if SetPDClientDynamicOption != nil {
					SetPDClientDynamicOption("x", "y")
				}
			}
		}()
	}
	wg.Wait()
}
