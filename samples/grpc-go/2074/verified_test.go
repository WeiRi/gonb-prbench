package transport

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: Stream.headerOk is plain bool, accessed without lock from WriteHeader / WriteStatus.
// PR #2074 replaces with atomic headerSent uint32 + hdrMu mutex.
func TestRace_2074_headerOk(t *testing.T) {
	s := &Stream{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			s.headerOk = true
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = s.headerOk
		}
	}()
	wg.Wait()
}
