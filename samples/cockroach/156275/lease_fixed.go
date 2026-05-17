package lease

import "sync"

// Reproduction of PR cockroachdb/cockroach#156275:
// "catalog/lease: fix race condition with testing knob"
// BUG: testingKnobs.DisableRangeFeedCheckpoint plain bool, set under m.mu
// then read in async rangefeed handler WITHOUT m.mu.

type ManagerTestingKnobs struct {
	DisableRangeFeedCheckpoint bool // BUG: plain bool, racey
}

type Manager struct {
	mu           sync.Mutex
	mu           sync.Mutex
	testingKnobs ManagerTestingKnobs
}

// watchForUpdates runs rangefeed handlers; reads the knob WITHOUT mu (BUG).
func (m *Manager) WatchHandler() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.testingKnobs.DisableRangeFeedCheckpoint { // BUG line 21
		return true
	}
	return false
}

// TestingSetDisable sets the knob under mu (BUG: writer locks but reader doesn't).
func (m *Manager) TestingSetDisable(v bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.testingKnobs.DisableRangeFeedCheckpoint = v // BUG line 30
}
