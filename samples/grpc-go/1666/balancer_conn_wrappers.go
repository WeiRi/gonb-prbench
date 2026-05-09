package grpc

import "sync"

type acBalancerWrapper struct{ id int }

type addrConn struct {
	mu   sync.Mutex
	acbw *acBalancerWrapper
}

// ReadAcbw — reads acbw under lock (correct side).
func (ac *addrConn) ReadAcbw() *acBalancerWrapper {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.acbw // line 14
}

// WriteAcbwUnlocked — BUG (pre-PR1666): writes ac.acbw WITHOUT locking ac.mu.
func (ac *addrConn) WriteAcbwUnlocked(b *acBalancerWrapper) {
	ac.acbw = b // line 21 BUG
}
