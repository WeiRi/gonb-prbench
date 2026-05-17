// Race test for traefik-2287 — HealthCheck.SetBackendsConfiguration race
package healthcheck

import (
	"context"
	"sync"
	"testing"
)

func TestRace_2287_SetBackends(t *testing.T) {
	hc := &HealthCheck{Backends: make(map[string]*BackendHealthCheck)}

	makeBackends := func() map[string]*BackendHealthCheck {
		return map[string]*BackendHealthCheck{"b1": nil, "b2": nil}
	}

	var wg sync.WaitGroup
	const N = 15
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 30; j++ {
				hc.SetBackendsConfiguration(context.Background(), makeBackends())
			}
		}()
	}
	wg.Wait()
}
