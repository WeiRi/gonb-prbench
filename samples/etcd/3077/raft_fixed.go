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

// FIX (PR3077): allocate stopped/done BEFORE spawning the run goroutine, so
// external readers see initialized channels.
func (s *EtcdServer) run() {
	s.r.stopped = make(chan struct{})
	s.r.done = make(chan struct{})
	go func() {
		close(s.r.stopped)
		close(s.r.done)
	}()
	_ = s.r.stopped
	_ = s.r.done
}
