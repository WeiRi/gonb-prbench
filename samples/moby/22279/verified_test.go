package container

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// BUG: WaitRunning reads s.waitChan after Unlock; if caller of SetRunning
// doesn't externally serialize with WaitRunning, race on s.waitChan close+
// reassign vs WaitRunning's wait(waitChan, timeout).
func TestRace_moby_22279_waitchan(t *testing.T) {
	s := NewState()
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
			// caller skipped Lock (mimics bug-prone caller pattern)
			s.SetRunning(j, false)
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
			_, _ = s.WaitRunning(time.Microsecond)
		}
	}()
	wg.Wait()
}
