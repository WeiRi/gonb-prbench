package etcdserver

import (
	"sync"
	"testing"
)

// TestRace_PR3077_RaftInitOrdering reproduces the data race fixed by
// https://github.com/etcd-io/etcd/pull/3077 — r.stopped / r.done are
// assigned inside r.run() goroutine while the spawning goroutine reads them.
func TestRace_PR3077_RaftInitOrdering(t *testing.T) {
	const N = 200
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			s := newEtcdServer()
			s.run()
		}()
	}
	wg.Wait()
}
