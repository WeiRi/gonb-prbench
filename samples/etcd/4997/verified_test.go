// Race test for etcd-4997 (etcdserver.consistentIndex)
// Targets the race fixed by adding atomic.Load/Store in consistent_index.go
// BUG: setConsistentIndex does `*i = consistentIndex(v)` (plain write),
//      ConsistentIndex does `return uint64(*i)` (plain read) — concurrent
//      access races on the uint64.
package etcdserver

import (
	"sync"
	"testing"
)

func TestRace_4997_ConsistentIndex(t *testing.T) {
	var ci consistentIndex
	const N = 50
	const ITERS = 1000
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				ci.setConsistentIndex(uint64(id*1000 + j))
			}
		}(i)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = ci.ConsistentIndex()
			}
		}()
	}
	wg.Wait()
}
