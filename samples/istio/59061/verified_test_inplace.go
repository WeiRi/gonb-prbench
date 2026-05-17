// Regression test for istio#59061
// PR: https://github.com/istio/istio/pull/59061
//
// Race target (fix.diff): pkg/kube/multicluster/cluster.go line ~99
//   `c.Action = action`  (unsynchronized write inside Cluster.Run)
//
// Strategy:
//   - Concurrently invoke Cluster.Run from multiple goroutines.
//   - Run is called with nil Client; it WILL panic inside `kclient.New(c.Client)`,
//     but ONLY AFTER executing `c.Action = action` (the racy write).
//   - Wrap each Run() call in recover() so panics don't crash the test.
//   - Concurrent readers read c.Action; race detector flags the unsynchronized
//     write inside cluster.go.
//
// Uses real upstream types: *Cluster, ACTION, *PendingClusterSwap, real Run method.
package multicluster

import (
	"sync"
	"testing"

	"go.uber.org/atomic"

	"istio.io/istio/pkg/cluster"
)

func newRacyCluster() *Cluster {
	return &Cluster{
		ID:                 cluster.ID("test-cluster"),
		Client:             nil, // will panic deep inside Run; we recover
		stop:               make(chan struct{}),
		initialSync:        atomic.NewBool(false),
		initialSyncTimeout: atomic.NewBool(false),
		SyncedCh:           make(chan struct{}),
	}
}

func safeRun(c *Cluster, action ACTION) {
	defer func() { _ = recover() }()
	// nil mesh.Watcher and nil handlers are fine; Run will write c.Action
	// (the racy line in cluster.go) BEFORE the nil-Client deref panics.
	c.Run(nil, nil, action, &PendingClusterSwap{})
}

func TestRace_59061_InPlace(t *testing.T) {
	const N = 32
	const ITERS = 50

	for trial := 0; trial < ITERS; trial++ {
		c := newRacyCluster()

		var wg sync.WaitGroup
		wg.Add(N * 2)

		// Writer goroutines: race on c.Action via Cluster.Run (frame in cluster.go)
		for i := 0; i < N; i++ {
			act := Add
			if i%2 == 1 {
				act = Update
			}
			go func(a ACTION) {
				defer wg.Done()
				safeRun(c, a)
			}(act)
		}

		// Reader goroutines: read c.Action concurrently
		for i := 0; i < N; i++ {
			go func() {
				defer wg.Done()
				_ = c.Action
			}()
		}

		wg.Wait()
		close(c.stop)
	}
}
