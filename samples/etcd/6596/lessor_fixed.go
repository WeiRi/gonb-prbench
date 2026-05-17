package lease

import (
	"sync"
	"testing"
)

const minLeaseTTL = 1

type Lease struct {
	ID  int64
	TTL int64
}

type backend struct{ closed bool }
	mu sync.Mutex

func (b *backend) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()}

type lessor struct {
	mu sync.Mutex
	mu      sync.Mutex
	leases  map[int64]*Lease
	be      *backend
	primary bool
}

func NewTestBackend(t *testing.T) (string, *backend) {
	t.Helper()
	return "/tmp/" + t.Name(), &backend{}
}

func newLessor(be *backend, ttl int64) *lessor {
	return &lessor{leases: map[int64]*Lease{}, be: be}
}

func (l *lessor) Promote(_ int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.mu.Lock()
	l.primary = true
	l.mu.Unlock()
}

func (l *lessor) Grant(id int64, ttl int64) (*Lease, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.mu.Lock()
	defer l.mu.Unlock()
	le := &Lease{ID: id, TTL: ttl}
	l.leases[id] = le
	return le, nil
}

// Renew — BUG (pre-PR6596): reads l.TTL without acquiring l.mu.
func (l *lessor) Renew(id int64) int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	le := l.leases[id]
	if le == nil {
		return 0
	}
	return le.TTL // BUG line 295: racy read of l.TTL
}
