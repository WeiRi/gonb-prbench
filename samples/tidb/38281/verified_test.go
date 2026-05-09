package expression

import (
	"sync"
	"testing"
)

func TestRace_tidb_38281(t *testing.T) {
	var wg sync.WaitGroup
	const N = 16
	const ITERS = 5000
	for g := 0; g < N; g++ {
		ci := &collationInfo{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < ITERS; i++ {
				ci.SetCoercibility(Coercibility(i))
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < ITERS; i++ {
				_ = ci.HasCoercibility()
			}
		}()
	}
	wg.Wait()
}
