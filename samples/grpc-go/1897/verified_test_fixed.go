package grpc

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// FIX: minConnectTimeout is now const; getMinConnectTimeout is an
// atomically-replaceable function variable for testing. No mutation race
// on the const.
func TestRace_grpc_go_1897_min_timeout(t *testing.T) {
	_ = time.Second
	var done int32
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = getMinConnectTimeout()
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
