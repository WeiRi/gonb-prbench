package handler

import (
	"sync"
	"testing"
)

// TestRace_4973 reproduces serving PR #4973 race in concurrency reporter:
// HandleRequest writes pendingRequests while ReportSnapshot iterates it
// without a mutex.
func TestRace_4973(t *testing.T) {
	for iter := 0; iter < 30; iter++ {
		r := NewConcurrencyReporter()

		var wg sync.WaitGroup
		// Writers
		for g := 0; g < 30; g++ {
			wg.Add(1)
			go func(gid int) {
				defer wg.Done()
				for i := 0; i < 100; i++ {
					r.HandleRequest("rev")
				}
			}(g)
		}
		// Readers
		for g := 0; g < 30; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 100; i++ {
					_ = r.ReportSnapshot()
				}
			}()
		}
		wg.Wait()
	}
}
