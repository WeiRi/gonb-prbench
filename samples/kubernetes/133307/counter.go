// Production stub for kubernetes-133307.
// Pre-PR: WithContext writes c.ctx, Add reads c.ctx via withExemplar — RACE.
package metrics

import (
	"context"
	"sync/atomic"

	"github.com/blang/semver/v4"
)

type CounterOpts struct {
	Namespace, Subsystem, Name, Help string
	ConstLabels                      map[string]string
}

type Counter struct {
	opts *CounterOpts
	ctx  context.Context // RACE: written by WithContext, read by Add via withExemplar
	val  atomic.Int64
}

func NewCounter(opts *CounterOpts) *Counter {
	return &Counter{opts: opts, ctx: context.Background()}
}

func (c *Counter) Create(_ *semver.Version) {}

// WithContext writes c.ctx (no lock) — RACE.
func (c *Counter) WithContext(ctx context.Context) *Counter {
	c.ctx = ctx
	return c
}

// Add calls withExemplar which reads c.ctx — RACE.
func (c *Counter) Add(v float64) {
	c.withExemplar()
	c.val.Add(int64(v))
}

func (c *Counter) withExemplar() {
	_ = c.ctx // RACE: read of ctx
}
