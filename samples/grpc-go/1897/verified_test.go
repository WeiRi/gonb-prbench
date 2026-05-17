package grpc

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// BUG: minConnectTimeout is a mutable package var; concurrent test code
// mutates it while production code reads. Race on the var.
func TestRace_grpc_go_1897_min_timeout(t *testing.T) {
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			minConnectTimeout = time.Duration(j) * time.Second
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = minConnectTimeout
		}
	}()
	wg.Wait()
}
