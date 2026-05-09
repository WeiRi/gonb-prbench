package net

import (
	"sync"
	"testing"
)

// TestRace_5865 reproduces serving PR #5865 race in activator revisionWatcher:
// CheckDests writes healthyPods while GetDests reads it without a mutex.
func TestRace_5865(t *testing.T) {
	for iter := 0; iter < 30; iter++ {
		r := NewRevisionWatcher()

		var wg sync.WaitGroup
		// Writers
		for g := 0; g < 30; g++ {
			wg.Add(1)
			go func(gid int) {
				defer wg.Done()
				for i := 0; i < 100; i++ {
					r.CheckDests("pod", true)
				}
			}(g)
		}
		// Readers
		for g := 0; g < 30; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 100; i++ {
					_ = r.GetDests()
				}
			}()
		}
		wg.Wait()
	}
}
