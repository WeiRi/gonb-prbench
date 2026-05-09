// Regression test for moby#21624 — volume/store/store.go data race
// PR: https://github.com/moby/moby/pull/21624
package main

import (
	"strconv"
	"sync"
	"testing"
)

func TestRace_21624(t *testing.T) {
	s := NewVolumeStore()
	const N = 50
	const ITERS = 100
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				s.create("vol-"+strconv.Itoa(id)+"-"+strconv.Itoa(j), map[string]string{"k": "v"})
			}
		}(i)
	}
	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = s.getVolume("vol-" + strconv.Itoa(id) + "-" + strconv.Itoa(j))
			}
		}(i)
	}
	wg.Wait()
}
