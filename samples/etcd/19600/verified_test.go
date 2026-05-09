package main

import (
	"sync"
	"testing"
)

// TestCancelCloseRace exercises concurrent create/cancel to trigger
// a data race on watcherGauge (plain int64).
// newWatcher() increments watcherGauge WITHOUT holding any lock.
// cancelWatcher() decrements watcherGauge UNDER s.mu.Lock().
// Concurrent Inc vs Dec = DATA RACE.
func TestCancelCloseRace(t *testing.T) {
	s := newWatchableStore()

	var wg sync.WaitGroup
	n := 50

	for g := 0; g < n; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				_, cancel := s.newWatcher(true)
				cancel() // calls cancelWatcher, which holds s.mu.Lock()
				// watcherGauge++ (in newWatcher, no lock)
				// races with watcherGauge-- (in cancelWatcher, under s.mu)
			}
		}()
	}

	wg.Wait()

	if watcherGauge != 0 {
		t.Logf("final watcherGauge = %d (expected 0)", watcherGauge)
	}
}
