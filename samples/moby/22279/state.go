// Production stub for moby container/state.go (PR #22279).
// Pre-PR: SetRunning closes & replaces s.waitChan WITHOUT s.Lock();
// WaitRunning reads s.waitChan under lock - race on field.
package container

import (
	"sync"
	"time"
)

type State struct {
	sync.Mutex
	running  bool
	pid      int
	waitChan chan struct{}
}

func NewState() *State {
	return &State{waitChan: make(chan struct{})}
}

// SetRunning is racy: writes pid + replaces waitChan without holding s.Lock().
func (s *State) SetRunning(pid int, initial bool) {
	s.running = true
	s.pid = pid
	close(s.waitChan)
	s.waitChan = make(chan struct{}) // RACE: write without lock
}

// WaitRunning reads s.waitChan under lock (race vs SetRunning's write).
func (s *State) WaitRunning(timeout time.Duration) error {
	s.Lock()
	wc := s.waitChan
	s.Unlock()
	select {
	case <-wc:
		return nil
	default:
		return nil
	}
}
