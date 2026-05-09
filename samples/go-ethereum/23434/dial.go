// PR #23434 - p2p/dial.go - data race on conn.flags field. Pre-fix: dial
// scheduler reads c.flags directly while another goroutine writes flags.
// PR fix: atomic.LoadInt32 / atomic.StoreInt32 on the field
// (and later commit removes the field from peers map entirely).
// Production-code path: p2p/dial.go (pre-fix line ~262).
package dial

import "sync"

type connFlag int32

type conn struct {
	flags connFlag
}

type dialScheduler struct {
	mu    sync.Mutex
	peers map[int]connFlag
}

func newScheduler() *dialScheduler {
	return &dialScheduler{peers: make(map[int]connFlag)}
}

// Pre-fix: register reads c.flags directly (no atomic) while concurrent
// goroutine sets flags via SetFlag.
// Upstream: p2p/dial.go (pre-fix line ~262).
func (d *dialScheduler) Register(id int, c *conn) {
	d.mu.Lock()
	d.peers[id] = c.flags // <- racy read of c.flags
	d.mu.Unlock()
}

// Pre-fix: SetFlag writes c.flags directly, no atomic.
func (c *conn) SetFlag(f connFlag) {
	c.flags = f
}
