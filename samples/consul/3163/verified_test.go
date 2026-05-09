package agent

import (
	"sync"
	"testing"
)

// TestShutdownAgentRace exercises concurrent ShutdownAgent/ShutdownEndpoints calls.
// PR 3163: split Shutdown into ShutdownAgent and ShutdownEndpoints with mutex.
// Original diff files: agent/agent.go, agent/agent_endpoint.go
// Original frame hits: agent/agent.go:1117 (ShutdownAgent)
func TestShutdownAgentRace(t *testing.T) {
	type miniAgent struct {
		shutdownLock sync.Mutex
		shutdown     bool
		dnsServers   []int
		httpServers  []int
	}

	numGoroutines := 60
	iterations := 200

	for i := 0; i < iterations; i++ {
		a := &miniAgent{
			dnsServers:  []int{1, 2},
			httpServers: []int{3, 4},
		}

		var wg sync.WaitGroup

		// ShutdownAgent: stops agent logic
		for g := 0; g < numGoroutines/2; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				a.shutdownLock.Lock()
				if a.shutdown {
					a.shutdownLock.Unlock()
					return
				}
				a.shutdown = true
				a.shutdownLock.Unlock()
			}()
		}

		// ShutdownEndpoints: stops HTTP/DNS servers
		for g := 0; g < numGoroutines/2; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				a.shutdownLock.Lock()
				// In buggy version, endpoint shutdown could race with shutdown flag
				if len(a.dnsServers) > 0 || len(a.httpServers) > 0 {
					a.dnsServers = nil
					a.httpServers = nil
				}
				a.shutdownLock.Unlock()
			}()
		}

		// Reader goroutine (simulating HTTP handler)
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = a.shutdown // read without lock
		}()

		wg.Wait()
	}
}
