package ui

import (
	"sync"
	"sync/atomic"
)

// Reproduction of PR cockroachdb/cockroach#149448 BUG state:
// spinner.go calls s.waitGroup.Add(1) INSIDE the spawned goroutine,
// after defer s.waitGroup.Done(). This races against waitGroup.Wait()
// that may be called before Add(1) executes.

type Spinner struct {
	waitGroup sync.WaitGroup
	flag      atomic.Int32
}

// Start: BUG — Add(1) is called inside the goroutine, after Done is deferred.
func (s *Spinner) Start() func() {
	go func() {
		defer s.waitGroup.Done()
		s.waitGroup.Add(1) // BUG: Add after Wait may have started
		s.flag.Store(1)
	}()
	return func() {
		s.waitGroup.Wait()
	}
}

