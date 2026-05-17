package server

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: c.traceMsg / c.traceOp do plain non-atomic read of package-level
// `trace int32`. Concurrent atomic.StoreInt32(&trace, ...) races with this read.
// PR #140 introduces per-client c.trace bool (initialized once in initClient),
// so subsequent reads check c.trace not package trace.
func TestRace_140_trace_var(t *testing.T) {
	c := &client{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		msg := []byte("MSG body")
		for i := 0; i < 10000 && atomic.LoadInt32(&done) == 0; i++ {
			c.traceMsg(msg)
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10000 && atomic.LoadInt32(&done) == 0; i++ {
			atomic.StoreInt32(&trace, int32(i&1))
		}
	}()
	wg.Wait()
	atomic.StoreInt32(&trace, 0)
}
