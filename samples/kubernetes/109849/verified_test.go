package v1beta1

import (
	"fmt"
	"sync"
	"testing"
)

// TestRace_109849_DisconnectMapRace reproduces the data race from PR #109849.
//
// BUG: disconnectClient reads s.clients[name] without holding the mutex,
// racing with registerClient/deregisterClient which write to the map.
func TestRace_109849_DisconnectMapRace(t *testing.T) {
	const N = 500
	const G = 8

	var wg sync.WaitGroup
	for g := 0; g < G; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < N; i++ {
				s := newServer()

				// Register some clients
				for j := 0; j < 10; j++ {
					name := fmt.Sprintf("plugin-%d", j)
					s.registerClient(name, &mockClient{name: name})
				}

				var inner sync.WaitGroup
				inner.Add(3)

				// Goroutine 1: continuously register and deregister clients
				go func() {
					defer inner.Done()
					for j := 0; j < 20; j++ {
						name := fmt.Sprintf("plugin-new-%d", j)
						s.registerClient(name, &mockClient{name: name})
						s.deregisterClient(name)
					}
				}()

				// Goroutine 2: call disconnectClient directly (unprotected map read)
				go func() {
					defer inner.Done()
					for j := 0; j < 20; j++ {
						s.disconnectClient("plugin-5")
					}
				}()

				// Goroutine 3: call DeRegisterPlugin (unprotected map read)
				go func() {
					defer inner.Done()
					for j := 0; j < 20; j++ {
						s.DeRegisterPlugin("plugin-3")
					}
				}()

				inner.Wait()
			}
		}()
	}
	wg.Wait()
}
