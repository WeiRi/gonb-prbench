// Production stub for nats-server server/client.go (PR #140).
// Pre-PR: traceOp reads global `trace` int32 non-atomically (line 210).
package server

var trace int32

type client struct{}

// traceOp non-atomically reads global trace.
func (c *client) traceOp(format, op string, arg []byte) {
	if trace != 0 { // RACE: bare read vs atomic.StoreInt32 in writers
		_ = format
		_ = op
		_ = arg
	}
}
