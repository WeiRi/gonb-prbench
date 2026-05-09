// Regression test for etcd#10772 — mvcc/index.go data race
// PR: https://github.com/etcd-io/etcd/pull/10772
package main

import (
	"sync"
	"testing"
)

func TestRace_10772(t *testing.T) {
	ti := newTreeIndex()
	ti.Put("racekey", 1)

	const N = 8
	const ITERS = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)

	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				ti.Put("racekey", int64(j+1))
			}
		}(i)
	}
	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = ti.Get("racekey")
			}
		}(i)
	}
	wg.Wait()
}
