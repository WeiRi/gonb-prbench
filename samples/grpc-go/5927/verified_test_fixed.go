package xdsclient

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX: ReportLoad keeps authorityMu held through refLocked. Caller-protected.
func TestRace_grpc_go_5927_refLocked_outside_lock(t *testing.T) {
	a := &authority{}
	var mu sync.Mutex
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			mu.Lock()
			a.refLocked()
			mu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			mu.Lock()
			_ = a.unrefLocked()
			mu.Unlock()
		}
	}()
	wg.Wait()
}
