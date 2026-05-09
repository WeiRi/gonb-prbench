// PR #23435 - miner/worker.go - data race between worker.close() and
// worker.mainLoop() on w.current.state. Pre-fix: close() (called by
// external goroutine) reads w.current and calls StopPrefetcher() while
// mainLoop() (own goroutine) reassigns w.current. PR fix: move
// StopPrefetcher to mainLoop()'s defer so w.current is touched by one
// goroutine only.
// Production-code path: miner/worker.go (pre-fix line ~321-323).
package miner

type stateDB struct {
	stopped bool
}

func (s *stateDB) StopPrefetcher() {
	s.stopped = true
}

type environment struct {
	state *stateDB
}

type worker struct {
	current *environment
	exitCh  chan struct{}
}

func newWorker() *worker {
	return &worker{exitCh: make(chan struct{}), current: &environment{state: &stateDB{}}}
}

// Close — pre-fix version: calls StopPrefetcher from external goroutine while
// mainLoop() reassigns w.current. Race on w.current.state pointer + the
// state.stopped field. Upstream: miner/worker.go (pre-fix line ~321).
func (w *worker) Close() {
	if w.current != nil && w.current.state != nil {
		w.current.state.StopPrefetcher()
	}
	close(w.exitCh)
}

// MainLoop — runs in own goroutine, periodically reassigns w.current.
// Upstream: miner/worker.go mainLoop() (pre-fix line ~446 area).
func (w *worker) MainLoop() {
	for {
		select {
		case <-w.exitCh:
			return
		default:
			w.current = &environment{state: &stateDB{}}
		}
	}
}
