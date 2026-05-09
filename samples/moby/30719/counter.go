// Production stub for moby daemon/graphdriver/counter.go (PR #30719).
// Pre-PR: Decrement reads counts without lock or doesn't acquire properly
// -> race vs Increment's writes.
package graphdriver

import "sync"

type Checker interface {
	IsMounted(path string) bool
}

type minfo struct {
	check bool
	count int
}

type RefCounter struct {
	mu      sync.Mutex
	counts  map[string]*minfo
	checker Checker
}

func NewRefCounter(c Checker) *RefCounter {
	return &RefCounter{counts: make(map[string]*minfo), checker: c}
}

// Increment writes counts[path] under lock.
func (c *RefCounter) Increment(path string) int {
	c.mu.Lock()
	m := c.counts[path]
	if m == nil {
		m = &minfo{}
		c.counts[path] = m
	}
	m.count++
	v := m.count
	c.mu.Unlock()
	return v
}

// Decrement reads counts[path] WITHOUT lock (pre-PR bug).
func (c *RefCounter) Decrement(path string) int {
	m := c.counts[path] // RACE: concurrent map read vs Increment writes
	if m == nil {
		return 0
	}
	m.count-- // RACE: write without lock
	return m.count
}
