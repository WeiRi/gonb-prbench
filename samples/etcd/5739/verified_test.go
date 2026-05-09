package etcdserver

import (
	"sync"
	"testing"
)

// TestRace_PR5739_AuthApplierUserUnlocked reproduces the data race fixed by
// https://github.com/etcd-io/etcd/pull/5739 — concurrent Apply() races on aa.user.
func TestRace_PR5739_AuthApplierUserUnlocked(t *testing.T) {
	aa := newAuthApplierV3()
	const N = 8
	const ITERS = 5000
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				aa.Apply("req", "u1")
			}
		}(i)
	}
	wg.Wait()
}
