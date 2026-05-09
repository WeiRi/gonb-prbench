package rafthttp

import (
	"sync"
	"testing"
)

// TestRace_PR2968_MsgAppReaderInit reproduces the data race fixed by
// https://github.com/etcd-io/etcd/pull/2968 — startPeer returns before
// p.msgAppReader is initialized in the goroutine, racing with external readers.
func TestRace_PR2968_MsgAppReaderInit(t *testing.T) {
	const N = 200
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			p := startPeer()
			_ = p.stop() // read p.msgAppReader concurrently with goroutine init
		}()
	}
	wg.Wait()
}
