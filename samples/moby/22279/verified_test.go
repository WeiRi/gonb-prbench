package container

import (
	"sync"
	"testing"
	"time"
)

func TestRace_22279(t *testing.T) {
	// Bug: SetRunning closes s.waitChan and creates a new one WITHOUT
	// holding s.Lock(). Concurrent WaitRunning calls read s.waitChan
	// under lock, creating a data race on s.waitChan field.
	s := NewState()

	var wg sync.WaitGroup
	nGoroutines := 50
	nIters := 200

	// Goroutines that call WaitRunning
	for i := 0; i < nGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				s.WaitRunning(-1 * time.Second)
			}
		}()
	}

	// Goroutines that call SetRunning (the write side of the race)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				s.SetRunning(j+100, false)
			}
		}()
	}

	wg.Wait()
}
