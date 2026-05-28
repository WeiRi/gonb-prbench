package apiserver

import (
	"sync"
	"testing"
)

// TestWarningWithRequestTimeout_115282: drive DefaultBuildHandlerChain().ServeRequest()
// concurrently. BUG: unsynchronized access to warningRecorder.warnings between
// the request goroutine (Add) and timeout goroutine (Snapshot).
func TestWarningWithRequestTimeout_115282(t *testing.T) {
	const N = 40
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := DefaultBuildHandlerChain()
			c.ServeRequest()
		}()
	}
	wg.Wait()
}
