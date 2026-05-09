package reconciler

import "sync"

// Stripped reproduction of pkg/kubelet/pluginmanager/reconciler/reconciler.go pre-PR #97193.
// BUG: AddHandler writes handlers map under sync.Mutex; getHandlers returns the same
// map without a copy, and callers iterate over it without lock => concurrent map access.

type reconciler struct {
	mu       sync.Mutex
	handlers map[string]interface{}
}

func (r *reconciler) AddHandler(name string, h interface{}) {
	r.mu.Lock()
	r.handlers[name] = h           // line 15 — locked write
	r.mu.Unlock()
}

// getHandlers — BUG: returns the underlying map, not a copy.
func (r *reconciler) getHandlers() map[string]interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.handlers              // line 23 — caller iterates without lock
}
