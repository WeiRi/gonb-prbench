package etcdserver

import "sync"

type EtcdServer struct {
	mu       sync.Mutex
	wg       sync.WaitGroup
	wgMu     sync.RWMutex
	stopping chan struct{}
	done     chan struct{}
}

func newEtcdServer() *EtcdServer {
	return &EtcdServer{
		stopping: make(chan struct{}),
		done:     make(chan struct{}),
	}
}

// run — BUG (pre-PR6628): wg.Wait() races with goAttach()'s wg.Add().
func (s *EtcdServer) run() {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer close(s.done)
	<-s.stopping
	s.wg.Wait() // line 35 BUG (no wgMu serialization)
}

// goAttach — BUG: wg.Add(1) without wgMu.
func (s *EtcdServer) goAttach(f func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-s.stopping:
		return
	default:
	}
	s.wg.Add(1) // line 39 BUG
	go func() {
		defer s.wg.Done()
		f() // line 43
	}()
}
