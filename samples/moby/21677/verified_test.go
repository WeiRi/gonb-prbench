// Regression test for moby#21677 — layer/ro_layer.go data race
// PR: https://github.com/moby/moby/pull/21677
package main

import (
	"sync"
	"testing"
)

func TestRace_21677(t *testing.T) {
	ls := NewLayerStore()
	l := ls.register("sha256:racy")
	const N = 50
	const ITERS = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				l.hold()
			}
		}()
	}
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = l.refCount()
			}
		}()
	}
	wg.Wait()
}
