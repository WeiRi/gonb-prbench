package transport

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX (PR #2074): headerOk replaced by atomic headerSent + hdrMu.
func TestRace_2074_headerOk(t *testing.T) {
	s := &Stream{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = s.updateHeaderSent()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = s.isHeaderSent()
		}
	}()
	wg.Wait()
}
