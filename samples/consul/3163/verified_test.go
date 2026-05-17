package agent

import (
	"sync"
	"testing"
)

func TestShutdownAgentRace(t *testing.T) {
	numGoroutines := 60
	iterations := 200

	for i := 0; i < iterations; i++ {
		a := &miniAgent{
			dnsServers:  []int{1, 2},
			httpServers: []int{3, 4},
		}

		var wg sync.WaitGroup

		for g := 0; g < numGoroutines/2; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				a.ShutdownAgent()
			}()
		}

		for g := 0; g < numGoroutines/2; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				a.ShutdownEndpoints()
			}()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = a.RacyReadShutdown()
		}()

		wg.Wait()
	}
}
