package ioutils

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

// Race: BUG Read does `bp.mu.Unlock(); return 0, bp.closeErr` reading
// closeErr after unlock; concurrent CloseWithError writes closeErr under
// lock. FIX captures err := bp.closeErr before Unlock.
func TestRace_moby_39445_bytespipe_closeerr(t *testing.T) {
	bp := NewBytesPipe()
	var done atomic.Bool
	var wg sync.WaitGroup
	buf := make([]byte, 4)

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000 && !done.Load(); i++ {
			_, _ = bp.Read(buf)
		}
	}()
	go func() {
		defer wg.Done()
		errs := []error{errors.New("e1"), errors.New("e2"), errors.New("e3")}
		for i := 0; i < 100000 && !done.Load(); i++ {
			_ = bp.CloseWithError(errs[i%3])
		}
		done.Store(true)
	}()
	wg.Wait()
}
