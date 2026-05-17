package grpc

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: line `ac.acbw = acbw` in UpdateAddresses else branch is unlocked.
// Concurrent readers under ac.mu race on ac.acbw field.
func TestRace_grpc_go_1666_acbw_write(t *testing.T) {
	ac := &addrConn{}
	acbw := &acBalancerWrapper{ac: ac}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			ac.acbw = acbw // BUG: no lock
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			ac.mu.Lock()
			_ = ac.acbw
			ac.mu.Unlock()
		}
	}()
	wg.Wait()
}
