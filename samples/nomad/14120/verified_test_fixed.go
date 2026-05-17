package nomad

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX (PR #14120): Bootstrapped is now Server.bootstrapped *atomic.Bool.
func TestRace_14120_bootstrapped(t *testing.T) {
	s := &Server{bootstrapped: &atomic.Bool{}}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			s.bootstrapped.Store(i&1 == 1)
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = s.bootstrapped.Load()
		}
	}()
	wg.Wait()
}
