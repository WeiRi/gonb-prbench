package wait

import (
	"sync"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/clock"
)

func TestRace_90476(t *testing.T) {
	// Use exponential backoff manager with real clock.
	// In BUG state, concurrent Backoff() calls race on:
	// 1. backoffTimer.Reset() - timer internal state
	// 2. getNextBackoff() - lastBackoffStart and backoff fields
	b := NewExponentialBackoffManager(
		1*time.Millisecond,
		1*time.Second,
		10*time.Second,
		2.0,
		0.0,
		clock.RealClock{},
	)

	const N = 30
	const ITERS = 200

	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = b.Backoff()
			}
		}()
	}

	wg.Wait()
}
