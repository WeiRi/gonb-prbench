package transport

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: do() has nested select default→{writes<-fn|closedCh}. Race window
// between outer select picking writes path and inner select sending — if
// another goroutine closes ht.writes (via WriteStatus' close(ht.writes)),
// panic: send on closed channel.
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
		defer func() {
			if r := recover(); r != nil {
				panic(r)
			}
		}()
		fn := func() {}
		for j := 0; j < 100 && atomic.LoadInt32(&done) == 0; j++ {
			_ = ht.do(fn)
		}
	}()
	go func() {
		defer wg.Done()
		close(ht.writes) // mimics WriteStatus's close
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
