package mvcc

import "sync"

type backend struct {
	mu sync.Mutex
}

func (b *backend) ForceCommit() {
	b.mu.Lock()
	b.mu.Unlock()
}

type store struct {
	mu sync.Mutex
	b  *backend
	hash uint64
}

func newStore() *store {
	return &store{b: &backend{}}
}

// Compact — calls ForceCommit under s.mu.
func (s *store) Compact(rev int64) {
	s.mu.Lock()
	s.hash = uint64(rev) // line 51 write under s.mu
	s.b.ForceCommit()
	s.mu.Unlock()
}

// Hash — BUG (pre-PR6279): reads s.hash and calls ForceCommit WITHOUT s.mu (line 57).
func (s *store) Hash() (uint64, error) {
	h := s.hash // BUG line 57: read without lock
	s.b.ForceCommit()
	return h, nil
}
