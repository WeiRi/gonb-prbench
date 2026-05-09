// Production stub for nats-server server/jetstream_cluster.go (PR #5182).
// Pre-PR: bare read of mset.lseq while only clMu held; mset.lseq is mset.mu-guarded.
package buggy

import "sync"

type stream struct {
	mu    sync.RWMutex
	clMu  sync.Mutex
	lseq  uint64
	clfs  uint64
}

// updateLseq writes lseq under mset.mu.
func (m *stream) updateLseq() {
	m.mu.Lock()
	m.lseq++
	m.mu.Unlock()
}

// processClusteredInboundMsg reads mset.lseq holding ONLY clMu (race vs updateLseq).
func (m *stream) processClusteredInboundMsg(_ int, _ int) {
	m.clMu.Lock()
	defer m.clMu.Unlock()
	_, _ = m.lseq, m.clfs // RACE: bare read of lseq (only clMu held)
}
