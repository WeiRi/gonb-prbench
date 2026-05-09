package fiber

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_PR4062_ValueAfterRelease(t *testing.T) {
	const iters = 200
	for i := 0; i < iters; i++ {
		c := NewDefaultCtx()

		var wg sync.WaitGroup
		var done int32

		// Goroutine A: simulates a leaked-handler goroutine that calls
		// ctx.Value() concurrently with the main request flow tearing down.
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { _ = recover() }()
			_ = c.Value("k")
			atomic.StoreInt32(&done, 1)
		}()

		// Goroutine B: simulates the main request flow finishing and the
		// pool reclaiming the ctx via app.ReleaseCtx() -> c.release().
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.release()
		}()

		wg.Wait()
		_ = atomic.LoadInt32(&done)
	}
}
