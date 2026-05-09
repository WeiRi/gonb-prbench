package net

// RevisionWatcher is a stand-in for knative/serving
// pkg/activator/net/revision_backends.go (BUG state).
// Bug: cb.healthyPods map is read by getDests and written by checkDests
// concurrently without synchronization.

type RevisionWatcher struct {
	healthyPods map[string]bool // BUG: unsynchronized
}

func NewRevisionWatcher() *RevisionWatcher {
	return &RevisionWatcher{healthyPods: map[string]bool{}}
}

// CheckDests writes healthyPods (race write).
func (r *RevisionWatcher) CheckDests(pod string, healthy bool) {
	r.healthyPods[pod] = healthy
}

// GetDests iterates healthyPods (race read).
func (r *RevisionWatcher) GetDests() int {
	count := 0
	for _, h := range r.healthyPods {
		if h {
			count++
		}
	}
	return count
}
