package scheduler

import "sync"

// Stripped reproduction of pkg/scheduler/schedule_one.go pre-PR #119729.
// BUG: handleBindingCycleError does NOT early-return; subsequent code calls
// Done(pi) which reads pi.Pod.UID, while a parallel "event handler" goroutine
// updates pi.Pod via writePod -> race.

type Pod struct {
	UID  string
	Name string
}

type PodInfo struct {
	Pod *Pod
}

type scheduler struct {
	wg sync.WaitGroup
}

// writePod simulates the eventHandler that updates pi.Pod after binding failure.
func (s *scheduler) writePod(pi *PodInfo) {        // line 25
	pi.Pod = &Pod{UID: "u2", Name: "n2"}             // line 26 — racing write
}

// done reads pi.Pod.UID — racing read.
func (s *scheduler) done(pi *PodInfo) string {     // line 30
	if pi.Pod == nil {
		return ""
	}
	return pi.Pod.UID                                // line 33 — racing read
}

// scheduleOne BUG body: missing early-return after handleBindingCycleError.
func (s *scheduler) scheduleOne(pi *PodInfo) {
	// emulate handleBindingCycleError dispatching async event handler
	s.wg.Add(1)
	go func() { defer s.wg.Done(); s.writePod(pi) }()
	// BUG: missing return; control falls through to Done()
	s.done(pi)
}
