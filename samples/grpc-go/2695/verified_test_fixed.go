package transport

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX: do() no longer has the inner select. close(ht.writes) is removed
// from WriteStatus. closedCh-only signaling means do() returns
// ErrConnClosing instead of panicking.
func TestRace_grpc_go_2695_send_closed(t *testing.T) {
	ht := &serverHandlerTransport{
		closedCh: make(chan struct{}),
		writes:   make(chan func(), 1),
	}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		fn := func() {}
		for j := 0; j < 100 && atomic.LoadInt32(&done) == 0; j++ {
			_ = ht.do(fn)
		}
	}()
	go func() {
		defer wg.Done()
		close(ht.closedCh)
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
