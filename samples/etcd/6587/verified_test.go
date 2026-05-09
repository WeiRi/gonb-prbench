package clientv3

import (
	"sync"
	"testing"
)

// TestRace_PR6587_InitReqRevRace reproduces the data race fixed by
// https://github.com/etcd-io/etcd/pull/6587 — serveSubstream defer writes
// ws.initReq.rev concurrently with resume() reading it.
func TestRace_PR6587_InitReqRevRace(t *testing.T) {
	const N = 30
	var wg sync.WaitGroup
	wg.Add(N * 3)
	for i := 0; i < N; i++ {
		w := newWatcherStream()
		go func() {
			defer wg.Done()
			w.serveSubstream() // drains respc, then defer writes initReq.rev
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				w.respc <- int64(j)
			}
			close(w.respc)
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				_ = w.resume()
			}
		}()
	}
	wg.Wait()
}
