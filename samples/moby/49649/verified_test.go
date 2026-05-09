package main

import (
	"sync"
	"testing"
)

func TestRace_49649(t *testing.T) {
	const N = 50
	const ITERS = 100

	tbl := newConnTrackTable()

	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				tbl.addConn(string(rune(j)), j)
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = tbl.lookupConn(string(rune(j)))
			}
		}()
	}
	wg.Wait()
}
