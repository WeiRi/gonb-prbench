// PR miekg/dns#1281 - client.go - data race on Client.Dialer.
// Pre-fix: ExchangeContext sets c.Dialer = &net.Dialer{Timeout:...} just
// before each call to Exchange, mutating shared *Client field. Concurrent
// ExchangeContext calls race write-vs-read on c.Dialer (Dial reads it).
// PR fix: plumb context.Context through DialContext/exchangeWithConnContext
// so deadlines are per-call without mutating shared Client state.
// Production-code path: client.go (pre-fix line ~440-449).
package dns

import (
	"net"
	"time"
)

type Client struct {
	Dialer *net.Dialer
}

// Dial — pre-fix path: reads c.Dialer.
// Upstream: client.go (pre-fix line ~85).
func (c *Client) Dial() *net.Dialer {
	if c.Dialer == nil {
		return &net.Dialer{}
	}
	return c.Dialer
}

// ExchangeContext — pre-fix: WRITES c.Dialer with a per-call timeout
// before invoking Exchange (which calls Dial). Concurrent ExchangeContext
// calls race on c.Dialer.
// Upstream: client.go (pre-fix line ~440-449).
func (c *Client) ExchangeContext(timeout time.Duration) *net.Dialer {
	c.Dialer = &net.Dialer{Timeout: timeout}
	return c.Dial()
}
