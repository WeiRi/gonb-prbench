package lease

import (
	"sync"
	"sync/atomic"
)

// Reproduction of PR cockroachdb/cockroach#144073:
// "catalog/lease: fix race condition releasing old versions".
// BUG: removeInactiveVersions reads/writes desc.mu.lease without taking desc.mu.

type storedLease struct{ id int }

type descriptorVersionState struct {
	refcount atomic.Int64
	mu       struct {
		sync.Mutex
		lease *storedLease
	}
}

type descriptorState struct {
	mu struct {
		sync.Mutex
		active []*descriptorVersionState
	}
}

// removeInactiveVersions accesses desc.mu.lease WITHOUT desc.mu (BUG).
func (t *descriptorState) removeInactiveVersions() []*storedLease {
	t.mu.Lock()
	defer t.mu.Unlock()
	var leases []*storedLease
	for _, desc := range t.mu.active {
		if desc.refcount.Load() == 0 {
			if l := desc.mu.lease; l != nil { // BUG line 32
				desc.mu.lease = nil           // BUG line 33
				leases = append(leases, l)
			}
		}
	}
	return leases
}

// SetLease is called concurrently by another path holding desc.mu.
func (d *descriptorVersionState) SetLease(l *storedLease) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.mu.lease = l
}

