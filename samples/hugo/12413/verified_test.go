// Race test for hugo PR #12413 — Wrap returning buf.Bytes() shares pooled buffer.
// Compatible with both BUG ([]byte) and FIX (string) return types via len()/range.
package hugocontext

import (
	"sync"
	"testing"
)

func TestRace_hugo_12413(t *testing.T) {
	const N = 64
	const ITERS = 500
	var wg sync.WaitGroup
	wg.Add(N)
	for g := 0; g < N; g++ {
		go func(id int) {
			defer wg.Done()
			payload := []byte("some markdown body content here for hugo wrap test")
			for i := 0; i < ITERS; i++ {
				out := Wrap(payload, uint64(id*1000+i))
				_ = len(out)
				var sum byte
				for k := 0; k < len(out); k++ {
					sum ^= out[k]
				}
				_ = sum
			}
		}(g)
	}
	wg.Wait()
}
