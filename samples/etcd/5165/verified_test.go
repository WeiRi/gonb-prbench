package raft

import (
	"sync"
	"testing"
)

func TestRace_5165(t *testing.T) {
	var wg sync.WaitGroup
	const numGoroutines = 50
	const iterations = 200

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				n := &node{
					propc:      make(chan Message, 1),
					recvc:      make(chan Message, 1),
					confc:      make(chan ConfChange, 1),
					confstatec: make(chan ConfState, 1),
					done:       make(chan struct{}),
					tickc:      make(chan struct{}, 1),
					status:     make(chan chan Status, 1),
					stop:       make(chan struct{}),
				}

				var inner sync.WaitGroup
				inner.Add(2)
				go func() {
					defer inner.Done()
					n.ClearPropc()
				}()
				go func() {
					defer inner.Done()
					_ = n.ReadPropc()
				}()
				inner.Wait()
			}
		}()
	}
	wg.Wait()
}
