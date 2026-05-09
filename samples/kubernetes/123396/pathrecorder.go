// Pre-fix pathrecorder.go from PR #123396 (apiserver/server/mux).
// BUG: ListedPaths reads m.exposedPaths WITHOUT m.lock while other goroutines
// may append; race on slice header.
package mux

import "sync"

type PathRecorderMux struct {
	lock         sync.Mutex
	exposedPaths []string
}

func NewPathRecorderMux() *PathRecorderMux {
	return &PathRecorderMux{}
}

// Handle / Register — writes m.exposedPaths under lock. pathrecorder.go:126
func (m *PathRecorderMux) Register(p string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.exposedPaths = append(m.exposedPaths, p) // line 126 — locked write
}

// ListedPaths — pathrecorder.go:99 pre-fix: NO lock on slice copy.
func (m *PathRecorderMux) ListedPaths() []string {
	handledPaths := append([]string{}, m.exposedPaths...) // line 99 — racy READ of slice header
	return handledPaths
}

// Helpers used in test, mirror pre-fix path lookup at :172/176/178.
func (m *PathRecorderMux) Has(p string) bool {
	for _, x := range m.exposedPaths { // line 172 — racy iter
		if x == p {
			return true
		}
	}
	return false
}
