package client

import (
	"sync"
	"time"

	"ase/etcd-20656/model"
)

type RecordingClient struct {
	watchMux        sync.Mutex
	watchOperations []model.WatchOperation
	baseTime        time.Time
}

// AppendRespUnsafe — BUG (pre-PR20656): appends to slice WITHOUT watchMux.
func (c *RecordingClient) AppendRespUnsafe(idx int, r model.WatchResponse) {
	c.watchOperations[idx].Responses = append(c.watchOperations[idx].Responses, r) // line 18 BUG
}

// ReadRespsUnsafe — BUG: iterates slice WITHOUT watchMux.
func (c *RecordingClient) ReadRespsUnsafe(idx int) int {
	n := 0
	for _, r := range c.watchOperations[idx].Responses { // line 24 BUG
		_ = r.Revision
		n++
	}
	return n
}

func (c *RecordingClient) InitWatch(req model.WatchRequest) {
	c.watchMux.Lock()
	c.watchOperations = append(c.watchOperations, model.WatchOperation{
		Request: req, Responses: []model.WatchResponse{},
	})
	c.watchMux.Unlock()
}
