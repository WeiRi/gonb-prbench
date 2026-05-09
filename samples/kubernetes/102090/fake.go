// Production stub for kubernetes-102090.
// Pre-PR: AddWatchReactor writes WatchReactionChain WITHOUT holding c.Lock().
package testing

import (
	"sync"

	"k8s.io/apimachinery/pkg/watch"
)

// Action is the minimal interface the test uses.
type Action interface{}

type WatchReactor interface {
	Handles(action Action) bool
	React(action Action) (bool, watch.Interface, error)
}

type WatchReactionFunc func(action Action) (handled bool, ret watch.Interface, err error)

type SimpleWatchReactor struct {
	Resource string
	Reaction WatchReactionFunc
}

func (s *SimpleWatchReactor) Handles(_ Action) bool { return true }
func (s *SimpleWatchReactor) React(action Action) (bool, watch.Interface, error) {
	return s.Reaction(action)
}

type Fake struct {
	sync.RWMutex
	WatchReactionChain []WatchReactor // RACE: read by InvokesWatch under lock; written by AddWatchReactor without lock
}

// AddWatchReactor (BUG): no Lock — RACE.
func (c *Fake) AddWatchReactor(resource string, reaction WatchReactionFunc) {
	c.WatchReactionChain = append(c.WatchReactionChain, &SimpleWatchReactor{Resource: resource, Reaction: reaction})
}
