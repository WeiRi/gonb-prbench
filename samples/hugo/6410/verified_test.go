// Race test for hugo-6410 — global DistinctLogger pointer race
// BUG: InitLoggers replaces DistinctErrorLog/WarnLog/FeedbackLog (pointer write);
//      concurrent .Println reads the pointer; race on global var
// FIX: InitLoggers calls .Reset() on existing logger; no pointer replacement
package helpers

import (
	"sync"
	"testing"
)

func TestRace_6410_DistinctLoggerInitRace(t *testing.T) {
	var wg sync.WaitGroup
	const N = 200
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			InitLoggers()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			DistinctErrorLog.Println("msg")
		}
	}()
	wg.Wait()
}
