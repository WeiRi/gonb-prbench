// Production stub for nats-server server/server.go (PR #6620).
// Pre-PR: s.sys.account read after s.mu.Unlock; fix snapshots sysAcc beforehand.
package buggy

import "sync"

type account struct {
	name string
}

type internal struct {
	account *account
}

type Server struct {
	mu  sync.RWMutex
	sys *internal
}

// rotateSysAccount writes s.sys.account under lock.
func (s *Server) rotateSysAccount(name string) {
	s.mu.Lock()
	s.sys = &internal{account: &account{name: name}}
	s.mu.Unlock()
}

// configureAccounts releases s.mu and THEN reads s.sys.account (pre-PR bug).
func (s *Server) configureAccounts() *account {
	s.mu.RLock()
	// some setup
	s.mu.RUnlock()
	// RACE: read s.sys.account after RUnlock
	if s.sys == nil {
		return nil
	}
	return s.sys.account
}
