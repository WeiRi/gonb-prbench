// Build tag for BUG state (bad is bool)


package pq

func newTestConn() *conn { return &conn{} }
func markBad(c *conn)    { c.bad = true }
func isBad(c *conn) bool { return c.bad }
