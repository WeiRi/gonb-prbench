// Regression test for moby#22966 — pkg/discovery/memory/memory.go data race
// PR: https://github.com/moby/moby/pull/22966
package main

import (
	"strconv"
	"sync"
	"testing"
)

func TestRace_22966(t *testing.T) {
	d := NewDiscovery()
	const N = 50
	const ITERS = 100
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				d.Register("host-" + strconv.Itoa(id) + "-" + strconv.Itoa(j))
			}
		}(i)
	}
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = d.Watch()
			}
		}()
	}
	wg.Wait()
}
