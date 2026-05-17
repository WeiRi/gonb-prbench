package server

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: populateGlobalPerSubjectInfo calls mb.readPerSubjectInfo(false) without
// holding mb.mu; FIX takes mb.mu.Lock. Race on mb internal state.
func TestRace_nats_3353_psim(t *testing.T) {
	mb := &msgBlock{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			mb.mu.Lock()
			mb.cache = nil
			mb.mu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = mb.cache
		}
	}()
	wg.Wait()
}
