package grpctest

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Race: BUG tLogger.ExpectErrorN, EndTest, expected access g.errors map
// without lock. Concurrent goroutines race on map writes.
// FIX adds g.m sync.Mutex around all map accesses.
func TestRace_grpc_go_3373_tlogger_errors(t *testing.T) {
	tl := TLogger
	tl.Update(t)

	var done int32
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200 && atomic.LoadInt32(&done) == 0; j++ {
				tl.ExpectErrorN("pat", 1)
			}
			atomic.StoreInt32(&done, 1)
		}()
	}
	wg.Wait()
}
