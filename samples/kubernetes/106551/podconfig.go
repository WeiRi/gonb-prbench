// Production stub for kubernetes-106551.
// Pre-PR: SeenAllSources reads c.sources without sourcesLock; Channel writes
// c.sources under sourcesLock — RACE.
package config

import (
	"context"
	"sync"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/record"
)

type PodConfigNotificationMode int

const (
	PodConfigNotificationIncremental PodConfigNotificationMode = iota
)

const TestSource = "test-source"

type PodConfig struct {
	mode         PodConfigNotificationMode
	recorder     record.EventRecorder
	sourcesLock  sync.Mutex
	sources      sets.String // RACE: written by Channel under lock, read by SeenAllSources without lock
}

func NewPodConfig(mode PodConfigNotificationMode, recorder record.EventRecorder) *PodConfig {
	return &PodConfig{mode: mode, recorder: recorder, sources: sets.NewString()}
}

// Channel writes c.sources under sourcesLock.
func (c *PodConfig) Channel(_ context.Context, src string) chan<- interface{} {
	c.sourcesLock.Lock()
	defer c.sourcesLock.Unlock()
	c.sources.Insert(src)
	return nil
}

// SeenAllSources reads c.sources WITHOUT sourcesLock — RACE.
func (c *PodConfig) SeenAllSources(seen sets.String) bool {
	for _, s := range c.sources.List() { // RACE
		if !seen.Has(s) {
			return false
		}
	}
	return true
}
