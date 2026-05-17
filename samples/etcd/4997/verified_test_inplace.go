// In-place race test for etcd-4997: package=etcdserver, uses upstream consistentIndex.
// Bug: consistent_index.go — setConsistentIndex writes without atomic,
// ConsistentIndex reads without atomic. Race on bare uint64.
// PR fix: use atomic.StoreUint64/LoadUint64.
package etcdserver

import (
	"sync"
	"testing"
)

func TestRace_4997_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 100

	var wg sync.WaitGroup
	wg.Add(N * 2)

	ci := new(consistentIndex)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				ci.setConsistentIndex(uint64(j)) // RACE WRITE (consistent_index.go:24)
			}
		}()
	}

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = ci.ConsistentIndex() // RACE READ (consistent_index.go:26)
			}
		}()
	}

	wg.Wait()
}
