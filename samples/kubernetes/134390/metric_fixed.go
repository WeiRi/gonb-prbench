package metrics

import (
	"sync"
	"sync/atomic"
)

type CounterOpts struct {
	Namespace, Subsystem, Name, Help string
}

// FIXED: isHidden/isDeprecated are atomic.Bool to avoid race with ClearState.
type Counter struct {
	createLock   sync.RWMutex
	isHidden     atomic.Bool
	isDeprecated atomic.Bool
}

func NewCounter(_ *CounterOpts) *Counter {
	return &Counter{}
}

func (c *Counter) ClearState() {
	c.createLock.Lock()
	c.isHidden.Store(!c.isHidden.Load())
	c.isDeprecated.Store(!c.isDeprecated.Load())
	c.createLock.Unlock()
}

func (c *Counter) IsHidden() bool {
	return c.isHidden.Load()
}

func (c *Counter) IsDeprecated() bool {
	return c.isDeprecated.Load()
}
