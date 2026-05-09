package pkg

import (
	"sync"
	"testing"
)

// TestRace_137949_UseStreamingBoolRace reproduces the data race from PR #137949.
//
// BUG: plain bool field read/written concurrently without synchronization.
// FIX: atomic.Bool with Load()/Store().
func TestRace_137949_UseStreamingBoolRace(t *testing.T) {
	const N = 5000
	const G = 10

	var wg sync.WaitGroup
	for g := 0; g < G; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc := &remoteImageService{useStreaming: true}

			var inner sync.WaitGroup
			inner.Add(2)

			// Goroutine 1: reads useStreaming
			go func() {
				defer inner.Done()
				for i := 0; i < N; i++ {
					_ = svc.ListImages()
				}
			}()

			// Goroutine 2: writes useStreaming = false
			go func() {
				defer inner.Done()
				for i := 0; i < N; i++ {
					svc.streamImagesFallback()
				}
			}()

			inner.Wait()
		}()
	}
	wg.Wait()
}
