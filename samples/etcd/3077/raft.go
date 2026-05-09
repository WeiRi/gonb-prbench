package etcdserver

import "sync"

type raftNode struct {
	stopped chan struct{}
	done    chan struct{}
	mu      sync.Mutex
}

type EtcdServer struct {
	r *raftNode
}

func newEtcdServer() *EtcdServer {
	return &EtcdServer{r: &raftNode{}}
}

// run — BUG (pre-PR3077): r.run goroutine assigns r.stopped/r.done while the
// caller already returned from start, racing with external readers.
func (s *EtcdServer) run() {
	go func() {
		// inside run goroutine — write r.stopped (line 26) and r.done (line 43)
		s.r.stopped = make(chan struct{})
		s.r.done = make(chan struct{})
		close(s.r.stopped)
		close(s.r.done)
	}()
	// concurrent reads of stopped/done from spawning goroutine
	_ = s.r.stopped
	_ = s.r.done
}
