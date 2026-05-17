package sql

import "sync"

// Reproduction of PR cockroachdb/cockroach#162460:
// "sql: fix a data race with parallel checks around metadata forwarder"
// BUG: createPhysPlan unconditionally writes planner.routineMetadataForwarder
// while parallel CHECK goroutines share the same planner.

type Forwarder struct{ id int }

var singletonNoop = &Forwarder{id: 1}

type Planner struct {
	routineMetadataForwarder *Forwarder // BUG: shared across parallel goroutines
}

type DistSQLPlanner struct{}
	mu sync.Mutex

// CreatePhysPlan (BUG): always mutates planner forwarder.
func (dsp *DistSQLPlanner) CreatePhysPlan(p *Planner) {
	dsp.mu.Lock()
	defer dsp.mu.Unlock()
	p.routineMetadataForwarder = singletonNoop // BUG line 20
	defer func() {
		p.routineMetadataForwarder = nil // BUG line 22
	}()
	_ = p.routineMetadataForwarder // line 24
}

// CheckRunner reads forwarder concurrently.
func CheckRunner(p *Planner) *Forwarder {
	return p.routineMetadataForwarder // BUG line 30
}

