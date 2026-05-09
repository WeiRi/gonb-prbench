package clientv3

import (
	"sync"

	"ase/etcd-13203/resolver"
)

type Config struct {
	Endpoints []string
}

type Client struct {
	mu       *sync.RWMutex
	cfg      Config
	resolver *resolver.Resolver
}

func (c *Client) SetEndpoints(eps ...string) {
	c.mu.Lock()
	c.cfg.Endpoints = append(c.cfg.Endpoints[:0], eps...)
	c.mu.Unlock()
}
