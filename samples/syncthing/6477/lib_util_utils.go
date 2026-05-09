package syncthing6477repro

import (
	"sync"
	"time"
)

type service struct {
	stopped chan struct{}
	mut     sync.Mutex
}

func newService() *service {
	s := &service{stopped: make(chan struct{})}
	close(s.stopped)
	return s
}

// lib/util/utils.go:235
func (s *service) Serve() {
	s.mut.Lock()
	s.stopped = make(chan struct{})
	stopped := s.stopped
	s.mut.Unlock()
	time.Sleep(time.Microsecond)
	s.mut.Lock()
	close(stopped)
	s.mut.Unlock()
}

// lib/util/utils.go:259 - reads s.stopped OUTSIDE lock
func (s *service) Stop() {
	s.mut.Lock()
	s.mut.Unlock()
	_ = <-s.stopped
}
