// VERIFIED race reproducer for nats-server PR #5182
// "Get last sequence with read lock in processClusteredInboundMsg"
// https://github.com/nats-io/nats-server/pull/5182
//
// Original racy file: server/jetstream_cluster.go (line 7703 pre-fix)
// Pre-fix: re-capture line `lseq, clfs = mset.lseq, mset.clfs` reads bare
// mset.lseq while only clMu is held (mset.lseq is guarded by mset.mu).
// Fix: replace `mset.lseq` with `mset.lastSeq()` which holds the RLock.
//
// 4 writers (lseq++) x 8 readers (clusteredInbound) x 5000 ops, golang:1.21 -race.
package buggy

import (
	"sync"
	"testing"
)

func TestRaceLseqRead(t *testing.T) {
	m := &stream{}
	const W = 4
	const R = 8
	const N = 5000
	var wg sync.WaitGroup
	for i := 0; i < W; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				m.updateLseq()
			}
		}()
	}
	for i := 0; i < R; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				m.processClusteredInboundMsg(0, 0)
			}
		}()
	}
	wg.Wait()
}
