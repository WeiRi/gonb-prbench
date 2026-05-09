package mvcc

import (
	"sync"
	"testing"
)

// TestRace_PR6279_ForceCommitOutsideLock reproduces the data race fixed by
// https://github.com/etcd-io/etcd/pull/6279 — Hash() calls b.ForceCommit
// without holding s.mu, racing with Compact() which calls ForceCommit under s.mu.
func TestRace_PR6279_ForceCommitOutsideLock(t *testing.T) {
	s := newStore()
	const N = 8
	const ITERS = 2000
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func(rev int64) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				s.Compact(rev + int64(j))
			}
		}(int64(i * ITERS))
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_, _ = s.Hash()
			}
		}()
	}
	wg.Wait()
}
