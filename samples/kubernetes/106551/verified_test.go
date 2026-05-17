// PR: https://github.com/kubernetes/kubernetes/pull/106551
// Fix: Add sourcesLock.Lock/Unlock in SeenAllSources to prevent data race
// on c.sources when called concurrently with Channel().
package config

import (
	"context"
	"strconv"
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	clientscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
)

func TestRace106551SeenAllSources(t *testing.T) {
	eventBroadcaster := record.NewBroadcaster()
	config := NewPodConfig(PodConfigNotificationIncremental, eventBroadcaster.NewRecorder(clientscheme.Scheme, v1.EventSource{Component: "kubelet"}))
	seenSources := sets.NewString(TestSource)

	var wg sync.WaitGroup
	const iterations = 100
	wg.Add(iterations * 2)

	for i := 0; i < iterations; i++ {
		go func(idx int) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			// Channel adds a source to c.sources, which SeenAllSources reads
			config.Channel(ctx, strconv.Itoa(idx))
		}(i)
		go func() {
			defer wg.Done()
			// SeenAllSources reads c.sources.List() which races with Channel
			// adding new sources in the bug
			config.SeenAllSources(seenSources)
		}()
	}

	wg.Wait()
}
