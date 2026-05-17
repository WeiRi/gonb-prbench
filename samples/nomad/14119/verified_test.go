package main

import (
	"sync"
	"testing"
	"time"
)

func TestRace_14119(t *testing.T) {
	const N = 55
	const ITERS = 250

	for trial := 0; trial < ITERS; trial++ {
		h := &heartbeatStop{
			shutdown: make(chan struct{}),
		}

		go h.watch()

		var wg sync.WaitGroup
		for i := 0; i < N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 50; j++ {
					h.setLastOk(time.Now())
					_ = h.lastOk // RACE READ without lock
				}
			}()
		}

		time.Sleep(2 * time.Millisecond)
		close(h.shutdown)
		wg.Wait()
	}
}
