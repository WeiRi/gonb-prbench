package v3rpc

import (
	"sync"
	"testing"
)

// TestRace_PR5897_ProgressPrevKVUnlocked reproduces the data race fixed by
// https://github.com/etcd-io/etcd/pull/5897 — recvLoop and sendLoop access
// progress / prevKV maps without holding sws.mu.
func TestRace_PR5897_ProgressPrevKVUnlocked(t *testing.T) {
	sws := newServerWatchStream()
	const N = 6
	const ITERS = 1000
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func(id WatchID) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				sws.recvLoop(id, true, true)
			}
		}(WatchID(i))
		go func(id WatchID) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = sws.sendLoop(id)
			}
		}(WatchID(i))
	}
	wg.Wait()
}
