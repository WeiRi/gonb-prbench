package server

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: debugSubscribers's deferred goroutine reads local var nsubs without
// atomic.Load while async callbacks atomic.AddInt32(&nsubs, n). Race on
// int32 var.
func TestRace_nats_1178_nsubs_atomic(t *testing.T) {
	var nsubs int32
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			atomic.AddInt32(&nsubs, 1)
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = nsubs // BUG: plain read
		}
	}()
	wg.Wait()
}
