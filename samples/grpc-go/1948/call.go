package grpc

type CallOption struct {
	Key string
	Val int
}

type dopts struct {
	callOptions []CallOption
}

type ClientConn struct {
	dopts dopts
}

// Invoke — BUG (pre-PR1948): append shares backing array when cap > len; concurrent
// callers race on the same slot.
func (cc *ClientConn) Invoke(opts ...CallOption) {
	merged := append(cc.dopts.callOptions, opts...) // line 29 BUG
	_ = merged
}
