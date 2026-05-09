package grpcproxy

import (
	"sync"
	"testing"
)

// TestRace_PR6704_GroupRevUnlocked reproduces the data race fixed by
// https://github.com/etcd-io/etcd/pull/6704 — maybeJoinWatcherSingle
// reads group.rev without group.mu, racing with broadcast() that writes
// group.rev under group.mu.
func TestRace_PR6704_GroupRevUnlocked(t *testing.T) {
	wgs := newWatcherGroups()
	g := newWatcherGroup()
	wgs.groups["k"] = g

	const N = 8
	const ITERS = 2000
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func(start int64) {
			defer wg.Done()
			for j := int64(0); j < ITERS; j++ {
				g.broadcast(start + j) // writes g.rev under g.mu
			}
		}(int64(i * ITERS))
		go func(streamID int64) {
			defer wg.Done()
			for j := int64(0); j < ITERS; j++ {
				wgs.maybeJoinWatcherSingle(
					receiverID{streamID: streamID, watcherID: j},
					watcherSingle{w: watcher{rev: j, wr: "k"}, sws: &sws{id: streamID}},
				)
			}
		}(int64(i))
	}
	wg.Wait()
}
