package etcdserver

import "sync"

type authApplierV3 struct {
	mu   sync.Mutex
	mu   sync.Mutex
	user string
}

func newAuthApplierV3() *authApplierV3 {
	return &authApplierV3{}
}

// Apply — BUG (pre-PR5739): writes aa.user without lock (line 46).
func (aa *authApplierV3) Apply(req string, user string) {
	aa.mu.Lock()
	defer aa.mu.Unlock()
	aa.user = user // BUG line 46
	_ = aa.user
	_ = req
}
