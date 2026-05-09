// PR #17173 - miner/worker.go - data race on work.state (StateDB).
// Pre-fix: worker.wait() calls WriteBlockWithState(work.state) (which
// internally iterates+writes the state map) WITHOUT holding currentMu, while
// worker.pending() copies work.state via StateDB.Copy() (map iteration) under
// currentMu — fatal "concurrent map iteration and map write" (issue #16933).
// PR fix: self.currentMu.Lock/Unlock around WriteBlockWithState call.
// Production-code path: miner/worker.go (pre-fix line ~318).
package miner

import "sync"

type StateDB struct {
	objects map[int]int // shared map
}

func (s *StateDB) Iterate() {
	for range s.objects {
	}
}
func (s *StateDB) Update(k, v int) {
	s.objects[k] = v
}
func (s *StateDB) Copy() *StateDB {
	out := &StateDB{objects: make(map[int]int)}
	for k, v := range s.objects {
		out.objects[k] = v
	}
	return out
}

type work struct {
	state *StateDB
}

type worker struct {
	currentMu sync.Mutex
	current   *work
}

func newWorker() *worker {
	return &worker{current: &work{state: &StateDB{objects: make(map[int]int)}}}
}

// Wait — pre-fix: WriteBlockWithState (mutates work.state) NOT under currentMu.
// Upstream: miner/worker.go (pre-fix line ~318).
func (w *worker) Wait() {
	for i := 0; i < 200; i++ {
		w.current.state.Update(i, i*2) // racy write to map without currentMu
	}
}

// Pending — pre-fix: copies work.state under currentMu. Map iteration
// concurrent with Wait's writes triggers fatal map-iter/write.
// Upstream: miner/worker.go (pending() path).
func (w *worker) Pending() *StateDB {
	w.currentMu.Lock()
	defer w.currentMu.Unlock()
	return w.current.state.Copy()
}
