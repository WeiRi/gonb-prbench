package transport

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Race: BUG inFlow.maybeAdjust path `if int32(n) > estSenderQuota`:
//   f.mu.Unlock(); return f.delta   → reads f.delta after unlock
// Concurrent maybeAdjust call writes f.delta = n under lock.
// FIX uses `defer f.mu.Unlock()` so f.delta read is inside lock.
func TestRace_grpc_go_2953_inflow_delta(t *testing.T) {
	// limit small enough that any n > 0 triggers the racy branch.
	f := &inFlow{limit: 64}

	var done int32
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
				_ = f.maybeAdjust(uint32(1 << 14))
			}
			atomic.StoreInt32(&done, 1)
		}()
	}
	wg.Wait()
}
