// Production stub for syncthing lib/nat/service.go (PR #8042).
// Pre-PR: process() reads Mapping fields under service mutex, but updateMapping
// modifies them without per-mapping synchronization.
package nat

import (
	"sync"
	"time"
)

type Address struct {
	IP   string
	Port int
}

type Mapping struct {
	mut          sync.RWMutex // (added by FIX, absent in BUG)
	extAddresses map[string]Address
	expires      time.Time
}

type Service struct {
	mut      sync.RWMutex
	mappings []*Mapping
}

// process reads mapping fields under service-level lock but no per-mapping lock.
func (s *Service) process() {
	s.mut.RLock()
	defer s.mut.RUnlock()
	for _, m := range s.mappings {
		// RACE: read m.expires and m.extAddresses without per-mapping lock
		_ = m.expires
		for k := range m.extAddresses {
			_ = k
		}
	}
}

// updateMapping writes mapping fields without holding service lock.
func (s *Service) updateMapping(m *Mapping) {
	m.expires = time.Now().Add(time.Hour) // RACE
	if m.extAddresses == nil {
		m.extAddresses = make(map[string]Address)
	}
	m.extAddresses["new"] = Address{IP: "1.2.3.4", Port: 1234} // RACE: concurrent map write
}
