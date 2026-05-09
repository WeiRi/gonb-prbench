// Production stub for consul agent/consul/client.go (PR #4216).
// Pre-PR: rpcLimiter *rate.Limiter is a plain pointer racing between
// RPC().Allow() readers and ReloadConfig() pointer writes.
package consul

// Limiter is a stand-in for rate.Limiter (avoids external module dep).
type Limiter struct {
	rate  int
	burst int
}

func NewLimiter(r, b int) *Limiter { return &Limiter{rate: r, burst: b} }
func (l *Limiter) Allow() bool      { return l.burst > 0 }

// Client mirrors the racy struct (line 60 in upstream client.go).
type Client struct {
	rpcLimiter *Limiter // racy: read in RPC, replaced in ReloadConfig
}

func NewClient() *Client {
	return &Client{rpcLimiter: NewLimiter(100, 1000)}
}

// RPC reads c.rpcLimiter without sync (line 263 upstream).
func (c *Client) RPC() bool {
	return c.rpcLimiter.Allow()
}

// ReloadConfig replaces the limiter pointer without sync.
func (c *Client) ReloadConfig() {
	c.rpcLimiter = NewLimiter(200, 2000)
}
