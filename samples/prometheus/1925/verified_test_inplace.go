// In-place race test for prometheus-1925: package=promql, uses upstream lexer.
// Bug: lex.go — lexer.seriesDesc (bare bool) written by caller after lex() starts
// goroutine, while l.run() reads it. PR fix adds atomic/sync for seriesDesc.
package promql

import (
	"sync"
	"testing"
)

func TestRace_1925_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 100
	for trial := 0; trial < ITERS; trial++ {
		l := lex("{} _ 1 x 5")
		var wg sync.WaitGroup
		for i := 0; i < N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				l.seriesDesc = true // RACE WRITE on bare bool (lex.go:341-343)
			}()
		}
		// Drain items until channel closed by l.run()
		for range l.items {
		}
		wg.Wait()
	}
}
