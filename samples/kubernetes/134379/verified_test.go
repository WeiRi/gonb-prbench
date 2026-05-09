package garbagecollector

import (
	"sync"
	"testing"
)

func TestRace_134379(t *testing.T) {
	n := &node{
		identity: objectReference{
			OwnerReference: OwnerReference{
				APIVersion: "v1", Kind: "Pod", Name: "test", UID: "test-uid",
			},
		},
		dependents: make(map[*node]struct{}),
	}

	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				n.String()
				n.markBeingDeleted()
			}
		}()
	}

	wg.Wait()
}
