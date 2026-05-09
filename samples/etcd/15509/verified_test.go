package embed

import (
	"sync"
	"testing"
)

// TestRace_PR15509_GsClosureCapture reproduces the data race fixed by
// https://github.com/etcd-io/etcd/pull/15509 — closure captures gs; two
// branches both assign gs and start a goroutine that reads it.
func TestRace_PR15509_GsClosureCapture(t *testing.T) {
	const N = 100
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go runServe(&wg)
	}
	wg.Wait()
}
