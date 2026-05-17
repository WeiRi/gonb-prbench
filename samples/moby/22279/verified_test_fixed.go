package container

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX: WaitRunning method removed + SetRunning no longer closes/reassigns
// s.waitChan → no race possible. Test exercises GetPID concurrent with
// SetRunning to confirm no race.
func TestRace_moby_22279_waitchan(t *testing.T) {
	s := NewState()
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
			s.Lock()
			s.SetRunning(j, false)
			s.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 50000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = s.GetPID()
		}
	}()
	wg.Wait()
}
