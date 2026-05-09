package etcdserver

import (
	"sync"
	"testing"
	"time"
)

// TestRace_PR6628_WaitGroupAddVsWait reproduces the sync.WaitGroup misuse
// fixed by https://github.com/etcd-io/etcd/pull/6628 — goAttach calls
// wg.Add(1) racing with run()'s wg.Wait().
//
// Race detector flags WaitGroup.Add(positive) racing with Wait directly via
// sync internal happens-before instrumentation, which is exactly what the
// PR was added to fix (the wgMu RWMutex serializes Add against close).
func TestRace_PR6628_WaitGroupAddVsWait(t *testing.T) {
	const N = 30
	for i := 0; i < N; i++ {
		s := newEtcdServer()
		runDone := make(chan struct{})
		go func() {
			defer func() { recover(); close(runDone) }()
			s.run()
		}()
		go func() {
			defer func() { recover() }()
			for j := 0; j < 200; j++ {
				s.goAttach(func() {})
			}
		}()
		time.Sleep(time.Microsecond)
		// race window has occurred during the Add+Wait overlap
		func() {
			defer func() { recover() }()
			close(s.stopping)
		}()
		select {
		case <-runDone:
		case <-time.After(500 * time.Millisecond):
		}
	}
	// Quiet the linter — variable use.
	var _ = sync.WaitGroup{}
}
