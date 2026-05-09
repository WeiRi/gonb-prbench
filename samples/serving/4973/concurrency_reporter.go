package handler

// concurrencyReporter is a stand-in for knative/serving
// pkg/activator/handler/concurrency_reporter.go (BUG state).
// Bug: pendingRequests map is read by ReportSnapshot and written by
// HandleRequest concurrently without synchronization.

type ConcurrencyReporter struct {
	pendingRequests map[string]int // BUG: unsynchronized
}

func NewConcurrencyReporter() *ConcurrencyReporter {
	return &ConcurrencyReporter{pendingRequests: map[string]int{}}
}

// HandleRequest writes to pendingRequests (race write).
func (r *ConcurrencyReporter) HandleRequest(rev string) {
	r.pendingRequests[rev] = r.pendingRequests[rev] + 1
}

// ReportSnapshot iterates pendingRequests (race read).
func (r *ConcurrencyReporter) ReportSnapshot() int {
	total := 0
	for _, v := range r.pendingRequests {
		total += v
	}
	return total
}
