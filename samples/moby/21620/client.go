// Production stub for moby libcontainerd/client.go (PR #21620).
// Pre-PR: lock() reads & writes containerMutexes without proper sync,
// causing concurrent map read/write race.
package libcontainerd

import "sync"

type container struct{}

type clientCommon struct {
	sync.Mutex
	backend          interface{}
	containers       map[string]*container
	containerMutexes map[string]*sync.Mutex
}

type client struct {
	clientCommon
}

// lock writes & reads the containerMutexes map without lock (pre-PR bug).
func (c *client) lock(id string) {
	if _, ok := c.containerMutexes[id]; !ok { // RACE: map read
		c.containerMutexes[id] = &sync.Mutex{} // RACE: map write
	}
	m := c.containerMutexes[id] // RACE: map read
	m.Lock()
}

func (c *client) unlock(id string) {
	m, ok := c.containerMutexes[id] // RACE: map read
	if ok {
		m.Unlock()
	}
}
