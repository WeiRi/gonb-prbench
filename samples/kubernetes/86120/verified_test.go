package watch

import (
	"sync"
	"testing"
)

func TestRace_86120(t *testing.T) {
	fw := NewFake()

	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				// Bug: Stopped is a public field read directly without lock,
				// while Stop/Reset write it under the lock.
				// fw.Stopped (direct read) races with fw.Stop() (locked write).
				_ = fw.Stopped
				fw.Stop()
				fw.Reset()
			}
		}()
	}

	wg.Wait()
}
