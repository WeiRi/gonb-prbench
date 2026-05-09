package xdsclient

import (
	"sync"

	"ase/grpc-go-5927/bootstrap"
	"ase/grpc-go-5927/cache"
)

type authority struct {
	refCount  int
	serverCfg *bootstrap.ServerConfig
}

type clientImpl struct {
	authorityMu     sync.Mutex
	authorities     map[string]*authority
	idleAuthorities *cache.TimeoutCache
}

// refLocked — increments a.refCount but BUG: caller already unlocked authorityMu.
func (c *clientImpl) refLocked(a *authority) {
	a.refCount++ // BUG: caller unlocked already
}

// ReportLoad — BUG (pre-PR5927): unlocks authorityMu BEFORE refLocked.
func (c *clientImpl) ReportLoad(cfg *bootstrap.ServerConfig) (interface{}, error) {
	c.authorityMu.Lock()
	a := c.authorities[cfg.String()]
	c.authorityMu.Unlock() // BUG: unlock before refLocked
	if a == nil {
		return nil, nil
	}
	c.refLocked(a)
	return nil, nil
}
