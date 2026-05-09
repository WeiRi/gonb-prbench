// Production stub for nomad client/client.go (PR #3461).
// Pre-PR: Client.Node() returns c.config.Node directly; updateAttributes
// rewrites Node.Attributes under configLock; readers iterate the same map
// after RLock release => race.
package buggy

import "sync"

type Node struct {
	ID         string
	Attributes map[string]string
}

type Config struct {
	Node *Node
}

type Client struct {
	configLock sync.RWMutex
	config     *Config
	configCopy *Config
}

// Node returns the racy shared *Node (pre-PR).
func (c *Client) Node() *Node {
	c.configLock.RLock()
	defer c.configLock.RUnlock()
	return c.config.Node
}

// updateAttributes writes Attributes map in place — racy vs map iteration.
func (c *Client) updateAttributes(k, v string) {
	c.configLock.Lock()
	c.config.Node.Attributes[k] = v
	c.configLock.Unlock()
}
