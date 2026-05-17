package vsphere

import "sync"

type Shared struct {
	mu  sync.Mutex
	val int64
}

func New() *Shared { return &Shared{} }
func (s *Shared) Write(v int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.val = v
}
func (s *Shared) Read() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.val
}
