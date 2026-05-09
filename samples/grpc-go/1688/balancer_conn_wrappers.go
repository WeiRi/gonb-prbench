package grpc

type acBalancerWrapper struct{ id int }

type ccBalancerWrapper struct {
	subConns map[*acBalancerWrapper]struct{}
}

// AddSC — BUG (pre-PR1688): writes subConns map without lock.
func (ccb *ccBalancerWrapper) AddSC(acbw *acBalancerWrapper) {
	ccb.subConns[acbw] = struct{}{} // line 12 BUG
}

// IterateSC — BUG: iterates subConns map without lock.
func (ccb *ccBalancerWrapper) IterateSC() int {
	n := 0
	for acbw := range ccb.subConns { // line 18 BUG
		_ = acbw
		n++
	}
	return n
}
