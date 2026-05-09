package libcontainerd

import (
	"sync"
	"testing"
)

func TestRace_21620(t *testing.T) {
	// Bug: lock() reads clnt.containerMutexes[containerID] on line 25
	// AFTER releasing clnt.Lock(). Another goroutine can concurrently
	// write to clnt.containerMutexes via lock()'s map assignment on line 22.
	// This causes a concurrent map read and map write, which is a data race
	// in Go (can corrupt the map and crash).
	clnt := &client{
		clientCommon: clientCommon{
			backend:          nil,
			containers:       make(map[string]*container),
			containerMutexes: make(map[string]*sync.Mutex),
		},
	}

	var wg sync.WaitGroup
	nGoroutines := 50
	nIters := 200

	for i := 0; i < nGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				cid := string(rune('a' + (id+j)%26))
				clnt.lock(cid)
				clnt.unlock(cid)
			}
		}(i)
	}

	wg.Wait()
}
