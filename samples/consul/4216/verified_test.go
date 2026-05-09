package consul

import (
	"sync"
	"testing"
)

func TestRPCLimiterRace(t *testing.T) {
	c := NewClient()
	numGoroutines := 60
	iterations := 200

	for i := 0; i < iterations; i++ {
		var wg sync.WaitGroup

		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = c.RPC()
			}()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			c.ReloadConfig()
		}()

		wg.Wait()
	}
}
