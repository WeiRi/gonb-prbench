package garbagecollector

import (
	"sync"
	"testing"
)

func TestRace_134379_NodeStringRace(t *testing.T) {
	n := &node{}
	var wg sync.WaitGroup
	const N = 200
	wg.Add(2)
	// Writer: set beingDeleted under lock
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			n.markBeingDeleted()
		}
	}()
	// Reader: String() reads fields including beingDeleted/virtual; BUG only RLocks dependentsLock
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			_ = n.String()
		}
	}()
	wg.Wait()
}
