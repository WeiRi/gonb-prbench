// In-place race test for prometheus-3738: package=file (discovery/file), uses upstream types.
// Bug: file.go — TimestampCollector.Collect iterates fileSD.timestamps without holding
// fileSD.lock, while writeTimestamp/deleteTimestamp modify timestamps under fileSD.lock.
// PR fix: add fileSD.lock.RLock() around iteration in Collect.
package file

import (
	"fmt"
	"sync"
	"testing"
)

func TestRace_3738_InPlace(t *testing.T) {
	d := &Discovery{
		timestamps: make(map[string]float64),
	}
	tc := &TimestampCollector{
		discoverers: map[*Discovery]struct{}{d: {}},
	}

	const N = 30
	var wg sync.WaitGroup
	wg.Add(N * 3)

	// Writers: modify timestamps under d.lock (via direct lock usage)
	for i := 0; i < N; i++ {
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				d.lock.Lock()
				d.timestamps[fmt.Sprintf("f%d", idx)] = float64(j)
				d.lock.Unlock()
			}
		}(i)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				d.lock.Lock()
				delete(d.timestamps, fmt.Sprintf("f%d", idx))
				d.lock.Unlock()
			}
		}(i)
		// Reader: calls Collect() which iterates timestamps WITHOUT fileSD.lock (BUG, file.go:98-99)
		go func() {
			defer wg.Done()
			out := make(map[string]float64)
			for j := 0; j < 50; j++ {
				tc.Collect(nil) // BUG: no lock on fileSD.timestamps (file.go:99)
			}
			_ = out
		}()
	}
	wg.Wait()
}
