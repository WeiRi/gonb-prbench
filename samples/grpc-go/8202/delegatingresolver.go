package delegatingresolver

import "sync"

type childResolver struct{ id int }

func (c *childResolver) ResolveNow() {}
func (c *childResolver) Close()      {}

type Resolver struct {
	mu             sync.Mutex
	targetResolver *childResolver
	proxyResolver  *childResolver
}

func New() *Resolver {
	r := &Resolver{}
	// Build path assigns fields without lock — BUG.
	r.targetResolver = &childResolver{id: 1} // line 33
	r.proxyResolver = &childResolver{id: 2}  // line 42
	return r
}

// updateProxyResolverState — calls into other child without childMu.
func (r *Resolver) updateProxyResolverState() {
	if r.targetResolver != nil {
		r.targetResolver.ResolveNow() // line 49 BUG: unlocked read
	}
}

// ResolveNow — calls both children unlocked.
func (r *Resolver) ResolveNow() {
	if r.proxyResolver != nil {
		r.proxyResolver.ResolveNow() // line 51
	}
	if r.targetResolver != nil {
		r.targetResolver.ResolveNow()
	}
}

// Close — BUG: nils both fields without lock.
func (r *Resolver) Close() {
	r.targetResolver = nil
	r.proxyResolver = nil
}
