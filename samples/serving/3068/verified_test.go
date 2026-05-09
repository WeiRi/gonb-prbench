package pocserving3068

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// Reproduces panic: send on closed channel from knative/serving pool.go (PR 3068).
// Oracle: PANIC. Frame target: pool.go (Go method).
func TestRace_serving3068(t *testing.T) {
	const iters = 200
	var panicked int32
	var firstStack string
	var stackMu sync.Mutex

	for it := 0; it < iters; it++ {
		p := NewWithCapacity(1, 5)
		var wg sync.WaitGroup
		// Concurrent producer
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.StoreInt32(&panicked, 1)
					buf := make([]byte, 4096)
					n := runtime.Stack(buf, false)
					stackMu.Lock()
					if firstStack == "" {
						firstStack = string(buf[:n])
					}
					stackMu.Unlock()
				}
			}()
			for j := 0; j < 50; j++ {
				p.Go(func() error { return nil })
			}
		}()
		// Concurrent Wait closes workCh
		go func() {
			p.Wait()
		}()
		wg.Wait()
		if atomic.LoadInt32(&panicked) == 1 {
			t.Logf("iter %d: panic stack:\n%s", it, firstStack)
			t.Fatalf("send on closed channel reproduced (pool.Go)")
		}
	}
}

