// Production stub for consul command/agent/gated_writer.go (PR #2262).
// GatedWriter.Write uses RLock when appending to w.buf — should be Lock.
package agent

import (
	"io"
	"sync"
)

type GatedWriter struct {
	Writer  io.Writer
	mu      sync.RWMutex
	buf     [][]byte
	flushed bool
}

// Write appends to buf under RLock — racy concurrent slice append (pre-PR bug).
func (gw *GatedWriter) Write(p []byte) (int, error) {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	if gw.flushed {
		return gw.Writer.Write(p)
	}
	p2 := make([]byte, len(p))
	copy(p2, p)
	gw.buf = append(gw.buf, p2) // RACE: concurrent append under RLock
	return len(p), nil
}
