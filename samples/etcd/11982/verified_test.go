package v2store

import (
	"sync"
	"testing"
)

// TestRace_11982 reproduces the race fixed by etcd#11982.
// : Lock domain incomplete (mixed atomic/plain access).
// BUG: clone() reads s.GetSuccess etc. as plain fields, while Inc()
//      writes them via atomic.AddUint64 -> race detector reports race.
// PR https://github.com/etcd-io/etcd/pull/11982 changes clone to use
// atomic.LoadUint64.
func TestRace_11982(t *testing.T) {
	s := newStats()
	const N = 30
	const ITERS = 200

	var wg sync.WaitGroup
	wg.Add(N * 2)

	// Writers: Inc (atomic.AddUint64 内部)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				s.Inc(SetSuccess)
				s.Inc(GetSuccess)
				s.Inc(DeleteSuccess)
			}
		}()
	}

	// Readers: clone 在 BUG 状态用 plain field read
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = s.clone()
			}
		}()
	}

	wg.Wait()
}
