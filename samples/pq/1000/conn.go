// Production stub for lib/pq conn.go (pre-PR #1000).
// Retains racy `bad bool` field; setBad/getBad model write/read sites.
package pq_1000_poc

type conn struct {
	bad bool
}

func (c *conn) setBad() { c.bad = true }

func (c *conn) getBad() bool { return c.bad }
