package kvserver

// Reproduction of PR cockroachdb/cockroach#67006:
// "kvserver: fix data race on replicaScanner.stopper".
// BUG: rs.stopper is assigned in Start() while concurrent code may read
// rs.stopper from a different goroutine (e.g. via stopper monitor).

type Stopper struct{ id int }

type replicaScanner struct {
	stopper *Stopper // BUG: written in Start without sync, also read elsewhere
}

// Start (BUG): writes rs.stopper inside Start while another goroutine
// may already be reading it.
func (rs *replicaScanner) Start(stopper *Stopper) {
	rs.stopper = stopper // BUG line 17
	rs.scanLoop()
}

// scanLoop runs in goroutine reading rs.stopper.
func (rs *replicaScanner) scanLoop() {
	if rs.stopper != nil { // line 23
		_ = rs.stopper.id
	}
}

// Monitor (BUG): peeks rs.stopper without sync.
func (rs *replicaScanner) Monitor() *Stopper {
	return rs.stopper // BUG line 30
}

