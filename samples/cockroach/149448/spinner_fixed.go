package ui

import (
	"sync"
	"sync/atomic"
)

type Spinner struct {
	waitGroup sync.WaitGroup
	flag      atomic.Int32
}

// FIX: Add(1) BEFORE go func, not inside.
func (s *Spinner) Start() func() {
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		s.flag.Store(1)
	}()
	return func() {
		s.waitGroup.Wait()
	}
}
