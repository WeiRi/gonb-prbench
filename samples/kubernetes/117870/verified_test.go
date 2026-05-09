package cache

import (
	"sync"
	"testing"
)

func TestRace_117870(t *testing.T) {
	s := &sharedIndexInformer{}

	var wg sync.WaitGroup
	n := 50

	// Readers: read s.transform and s.watchErrorHandler WITHOUT lock
	// (simulates Run() reading these fields before acquiring startedLock)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				_ = s.transform
				_ = s.watchErrorHandler
			}
		}()
	}

	// Writers: use SetTransform and SetWatchErrorHandler which
	// write under startedLock - but readers don't hold the lock
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				s.SetTransform(func(interface{}) (interface{}, error) {
					return nil, nil
				})
				s.SetWatchErrorHandler(func(r *Reflector, err error) {})
			}
		}()
	}

	wg.Wait()
}
