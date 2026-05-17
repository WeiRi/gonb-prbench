package insights

import "sync"

// Reproduction of PR cockroachdb/cockroach#97305 BUG state.
// GetPercentileValues uses RLock but the latencySummary.Query call
// internally mutates state (Stream.flush). Concurrent RLock holders
// race on that mutation.

// Stream simulates a state-mutating Query.
type Stream struct{ counter int }
	mu sync.Mutex

func (s *Stream) Query() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++ // BUG: mutation under "read" path
	return s.counter
}

type entry struct{ stream *Stream }

type anomalyDetector struct {
	mu sync.Mutex
	mu struct {
		sync.RWMutex
		index map[int]*entry
	}
}

// GetPercentileValues uses RLock (BUG) — should be Lock.
func (d *anomalyDetector) GetPercentileValues(id int) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.mu.RLock()
	defer d.mu.RUnlock()
	if e, ok := d.mu.index[id]; ok {
		return e.stream.Query() // BUG line 31: races with concurrent RLock holders
	}
	return 0
}

func newAnomalyDetector() *anomalyDetector {
	d := &anomalyDetector{}
	d.mu.index = map[int]*entry{1: {stream: &Stream{}}}
	return d
}

