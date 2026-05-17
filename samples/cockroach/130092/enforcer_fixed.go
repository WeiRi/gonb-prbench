package license

import (
	"sync"
)

type Enforcer struct {
	mu    sync.Mutex
	db    interface{}
	notes string
}

var instance *Enforcer
var once sync.Once

// Fix: use sync.Once.Do without prior unsynchronized read of instance.
func GetEnforcerInstance() *Enforcer {
	once.Do(func() {
		instance = newEnforcer()
	})
	return instance
}

func newEnforcer() *Enforcer {
	return &Enforcer{notes: "x"}
}

func (e *Enforcer) Start(db interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.db = db
}

func (e *Enforcer) DB() interface{} {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.db
}
