package slinstance

import "sync"

// Reproduction of PR cockroachdb/cockroach#69290:
// "sqlliveness/slinstance: fix data race"
// BUG: session.exp is at struct level, written by extendSession under l.mu
// but read by Expiration() WITHOUT lock.

type session struct {
	id  int
	exp int64 // BUG: not under any lock from session perspective
	mu  struct {
		sync.RWMutex
		callbacks []func()
	}
}

// Expiration reads exp WITHOUT mu (BUG).
func (s *session) Expiration() int64 {
	return s.exp // BUG line 21
}

type Instance struct {
	mu struct {
		sync.Mutex
		s *session
	}
}

// extendSession writes exp under l.mu (BUG: writer locks but reader doesn't).
func (l *Instance) ExtendSession(s *session, exp int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.mu.s.exp = exp // BUG line 32
}

