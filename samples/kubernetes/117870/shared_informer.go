package cache

import "sync"

// Stripped reproduction of shared_informer.go pre-PR #117870.
// BUG: SetTransform/SetWatchErrorHandler write under startedLock; Run() reads them WITHOUT.
// (The test reads them directly, simulating Run's unprotected read.)

type Reflector struct{}

type sharedIndexInformer struct {
	startedLock        sync.Mutex
	transform          func(interface{}) (interface{}, error)
	watchErrorHandler  func(*Reflector, error)
}

func (s *sharedIndexInformer) SetTransform(t func(interface{}) (interface{}, error)) {
	s.startedLock.Lock()
	s.transform = t
	s.startedLock.Unlock()
}

func (s *sharedIndexInformer) SetWatchErrorHandler(h func(*Reflector, error)) {
	s.startedLock.Lock()
	s.watchErrorHandler = h
	s.startedLock.Unlock()
}
