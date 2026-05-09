package metrics

import "sync"

// Stripped reproduction of staging/src/k8s.io/component-base/metrics/metric.go pre-PR #134390.
// BUG: ClearState() writes isHidden/isDeprecated under createLock; IsHidden()/IsDeprecated() read them WITHOUT the lock.

type CounterOpts struct {
	Namespace, Subsystem, Name, Help string
}

type Counter struct {
	createLock   sync.RWMutex
	isHidden     bool
	isDeprecated bool
}

func NewCounter(_ *CounterOpts) *Counter {
	return &Counter{}
}

func (c *Counter) ClearState() {
	c.createLock.Lock()
	c.isHidden = !c.isHidden
	c.isDeprecated = !c.isDeprecated
	c.createLock.Unlock()
}

// Plain (un-locked) read — BUG.
func (c *Counter) IsHidden() bool {
	return c.isHidden
}

func (c *Counter) IsDeprecated() bool {
	return c.isDeprecated
}
