package util

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	stsync "github.com/syncthing/syncthing/lib/sync"
)

// FIX: Stop() caches stopped under mut.Lock then unlocks and uses cached.
func TestRace_syncthing_6477_stopped_field(t *testing.T) {
	s := &service{
		ctx:     context.Background(),
		stopped: make(chan struct{}),
		mut:     stsync.NewMutex(),
	}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			s.mut.Lock()
			stopped := s.stopped
			s.mut.Unlock()
			_ = stopped
		}
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			s.mut.Lock()
			s.stopped = make(chan struct{})
			s.mut.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
