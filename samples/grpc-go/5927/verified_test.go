package xdsclient

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: ReportLoad releases authorityMu BEFORE calling a.refLocked() →
// race on authority.refCount.
func TestRace_grpc_go_5927_refLocked_outside_lock(t *testing.T) {
	a := &authority{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			a.refLocked()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = a.unrefLocked()
		}
	}()
	wg.Wait()
}
