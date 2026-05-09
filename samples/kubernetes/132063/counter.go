// Production stub for kubernetes-132063.
// Pre-PR: LabelValueAllowLists is read at counter.go:215 and written inside
// sync.Once at counter.go:221, racing across goroutines.
package metrics

import (
	"sync"

	"github.com/blang/semver/v4"
)

type allowListMap = map[string]string

var globalAllowList allowListMap

func SetLabelAllowListFromCLI(m map[string]string) { globalAllowList = m }

type CounterOpts struct {
	Namespace, Subsystem, Name, Help string
	ConstLabels                      map[string]string
}

type CounterVec struct {
	opts                  *CounterOpts
	labels                []string
	once                  sync.Once
	LabelValueAllowLists  allowListMap // RACE: read outside Do, written inside
}

func NewCounterVec(opts *CounterOpts, labels []string) *CounterVec {
	return &CounterVec{opts: opts, labels: labels}
}

func (c *CounterVec) Create(_ *semver.Version) {}

func (c *CounterVec) WithLabelValues(_ ...string) *CounterVec {
	if c.LabelValueAllowLists != nil { // RACE: read of map header
		_ = c.LabelValueAllowLists
	}
	c.once.Do(func() {
		c.LabelValueAllowLists = globalAllowList // RACE: write
	})
	return c
}
