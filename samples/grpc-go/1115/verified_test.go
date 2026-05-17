package transport

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: serverHandlerTransport.do() select picks `ht.writes <- fn` even when
// closedCh is just-closed if both cases are ready. After the chosen write
// proceeds, if another goroutine closes ht.writes too, panic: send on closed
// channel.
//
// Reproduce: Write to ht.writes concurrently while closing ht.writes from
// another goroutine — mimics the production race window.
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
		defer func() {
			if r := recover(); r != nil {
				panic(r) // re-raise to bubble up
			}
		}()
		fn := func() {}
		for j := 0; j < 100 && atomic.LoadInt32(&done) == 0; j++ {
			_ = ht.do(fn)
		}
	}()
	go func() {
		defer wg.Done()
		close(ht.closedCh)
		close(ht.writes)
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
