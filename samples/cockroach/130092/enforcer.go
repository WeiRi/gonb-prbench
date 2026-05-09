package license

import (
	"sync"
)

// Reproduction of PR cockroachdb/cockroach#130092 BUG state
// (sql/license: Fix data race in license/enforcer.go).
// Pre-fix code reads `instance` outside of once.Do (race) and writes
// e.db inside Start without synchronization (race vs concurrent uses
// of the singleton).

type Enforcer struct {
	db    interface{} // BUG: written in Start without sync
	notes string
}

var instance *Enforcer
var once sync.Once

// GetEnforcerInstance returns singleton (BUG state).
func GetEnforcerInstance() *Enforcer {
	if instance == nil { // BUG: read race vs once.Do writer
		once.Do(func() {
			instance = newEnforcer()
		})
	}
	return instance
}

func newEnforcer() *Enforcer {
	return &Enforcer{notes: "x"}
}

// Start writes e.db without sync (BUG).
func (e *Enforcer) Start(db interface{}) {
	e.db = db
}

// DB reads e.db without sync (BUG).
func (e *Enforcer) DB() interface{} {
	return e.db
}

