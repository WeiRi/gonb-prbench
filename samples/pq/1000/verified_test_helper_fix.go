// Build tag for FIX state (bad is *atomic.Value)
//go:build pqfix

package pq

import "sync/atomic"

func newTestConn() *conn { v := &atomic.Value{}; v.Store(false); return &conn{bad: v} }
func markBad(c *conn)    { c.bad.Store(true) }
func isBad(c *conn) bool { return c.bad.Load().(bool) }
