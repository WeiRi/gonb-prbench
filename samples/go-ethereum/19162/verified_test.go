package pss

import (
	"sync"
	"testing"
)

// TestRace_19162_HandshakeController_symKeyIndex: handler reads map without
// ctl.lock while ReleaseKey deletes under ctl.lock — concurrent map read+write.
func TestRace_19162_HandshakeController_symKeyIndex(t *testing.T) {
	ctl := NewHandshakeController()
	for i := 0; i < 8; i++ {
		ctl.Insert("k")
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			ctl.Insert("k")
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			ctl.Handler("k")
		}
	}()
	wg.Wait()
}
