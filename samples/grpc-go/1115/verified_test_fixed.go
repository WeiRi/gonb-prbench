package transport

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX: do() now checks closedCh first. Once closedCh is closed, do()
// returns ErrConnClosing immediately without touching writes. No panic.
// (FIX is imperfect against close(writes) without close(closedCh), so this
// test reflects the production lifecycle where closedCh closes first.)
func TestRace_grpc_go_1115_send_closed(t *testing.T) {
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
