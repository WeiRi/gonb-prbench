package pocgrpc2695

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_grpcgo2695(t *testing.T) {
	const iters = 400
	var panicked int32
	var firstStack string
	var stackMu sync.Mutex

	for it := 0; it < iters; it++ {
		ht := newServerHandlerTransport()

		runDone := make(chan struct{})
		go func() {
			ht.runStream()
			close(runDone)
		}()

		var wg sync.WaitGroup
		for w := 0; w < 6; w++ {
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
				for i := 0; i < 50; i++ {
					_ = ht.do(func() {})
				}
			}()
		}

		go func() {
			ht.WriteStatus()
		}()

		wg.Wait()
		<-runDone
		if atomic.LoadInt32(&panicked) == 1 {
			t.Logf("iter %d: caught panic stack:\n%s", it, firstStack)
			t.Fatalf("send on closed channel reproduced")
		}
	}
}
