package xdsclient

import "sync"

type lrsClient struct{ id int }

type clientImpl struct {
	mu        sync.Mutex
	lrsClient *lrsClient
}

func NewClientImpl() *clientImpl { return &clientImpl{} }

// ReportLoad — BUG (pre-PR8540/8541): lazy-init c.lrsClient WITHOUT lock.
// Concurrent callers race on the field assignment.
func (c *clientImpl) ReportLoad() *lrsClient {
	if c.lrsClient == nil { // line 20 BUG: racy read
		c.lrsClient = &lrsClient{id: 1} // line 25/27 BUG: racy write
	}
	return c.lrsClient
}
