package ioutils

import (
	"sync"
	"testing"
)

func TestRace_39445(t *testing.T) {
	const TRIALS = 100
	const READERS = 50
	for trial := 0; trial < TRIALS; trial++ {
		bp := NewBytesPipe()
		var started sync.WaitGroup
		var wg sync.WaitGroup

		// Start all readers first; they will block on wait condition
		started.Add(READERS)
		wg.Add(READERS)
		for i := 0; i < READERS; i++ {
			go func() {
				started.Done() // signal we've been spawned
				defer wg.Done()
				buf := make([]byte, 64)
				bp.Read(buf)
			}()
		}

		// Wait for all readers to be spawned before closing, so they
		// will all be either in the Read lock path or waiting on the
		// condition variable.
		started.Wait()

		// Now close the pipe - this wakes all waiting readers and
		// creates the race window on bp.closeErr at bytespipe.go:132
		bp.Close()

		// Some readers may have missed the broadcast if they didn't
		// yet reach the Wait() call. Close again (no-op in pre-fix)
		// to ensure they wake up too. Also call concurrently with
		// exiting readers to widen the race window.
		wg.Add(1)
		go func() {
			defer wg.Done()
			bp.Close()
		}()

		wg.Wait()
	}
}
