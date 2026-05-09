package tcpproxy

import "sync"

type remote struct {
	addr   string
	id     int
	active bool
}

func (r *remote) isActive() bool          { return r.active }
func (r *remote) tryReactivate() error    { _ = r.addr; _ = r.id; return nil }

type TCPProxy struct {
	mu      sync.Mutex
	remotes []*remote
}

// runMonitorOnce — BUG (pre-PR7361): closure captures range var r (line 33-34).
func (tp *TCPProxy) runMonitorOnce() {
	tp.mu.Lock()
	for _, r := range tp.remotes { // line 31
		if !r.isActive() {
			go func() { // line 33: captures r
				_ = r.tryReactivate() // line 34
			}()
		}
	}
	tp.mu.Unlock()
}
