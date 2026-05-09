package client

import (
	"sync"
	"testing"
)

// TestRace_47961 reproduces moby PR #47961 race in client/client.go:
// customHTTPHeaders map is concurrently mutated by SetCustomHTTPHeaders
// while reads happen on every API call, without synchronization.
func TestRace_47961(t *testing.T) {
	for iter := 0; iter < 30; iter++ {
		c := NewClient()

		var wg sync.WaitGroup
		// Writers
		for g := 0; g < 30; g++ {
			wg.Add(1)
			go func(gid int) {
				defer wg.Done()
				for i := 0; i < 100; i++ {
					c.AddHeader("k", "v")
				}
			}(g)
		}
		// Readers
		for g := 0; g < 30; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 100; i++ {
					_ = c.CountHeaders()
				}
			}()
		}
		wg.Wait()
	}
}
