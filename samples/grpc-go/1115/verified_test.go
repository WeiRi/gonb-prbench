// Race-trigger test for grpc-go-1115; see README.md for usage.

package transport

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_PR1115(t *testing.T) {
	const iters = 200
	for i := 0; i < iters; i++ {
		ht := &serverHandlerTransport{
			closedCh: make(chan struct{}),
			writes:   make(chan func()),
		}

		// Pre-close closedCh: simulates Close() already having been invoked
		// from the client-disconnect path.
		close(ht.closedCh)

		var wg sync.WaitGroup
		// Goroutine A: simulate the tail of a completed WriteStatus that runs
		// close(ht.writes) right at the moment another do() races in.
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { _ = recover() }()
			close(ht.writes)
		}()

		// Goroutine B: a racing Write/WriteHeader calling do() — BUG hits
		// 'send on closed channel' from handler_server.go's do().
		wg.Add(1)
		var bDone int32
		go func() {
			defer wg.Done()
			defer func() {
				_ = recover() // swallow runtime panic from BUG send-on-closed
				atomic.StoreInt32(&bDone, 1)
			}()
			_ = ht.do(func() {})
		}()

		wg.Wait()
		_ = atomic.LoadInt32(&bDone)
	}
}
