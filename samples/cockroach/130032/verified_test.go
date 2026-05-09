package diagnostics

import (
	"sync"
	"testing"
)

func Test130031Race(t *testing.T) {
	r := &Reporter{}
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				r.ReportDiagnostics()
			} else {
				_ = r.Read()
			}
		}(i)
	}
	wg.Wait()
}
