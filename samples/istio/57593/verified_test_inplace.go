// Regression test for istio#57593
// PR: https://github.com/istio/istio/pull/57593
// Race: sanitizeServerHostNamespace modifies server.Hosts slice in place
//        (at BUG state, no clone). Concurrent calls race on Hosts[i] writes.
// Uses real upstream types: networking.Server, real upstream method: sanitizeServerHostNamespace.
package model

import (
	"sync"
	"testing"

	networking "istio.io/api/networking/v1alpha3"
)

func TestRace_57593_InPlace(t *testing.T) {
	const N = 32
	const ITERS = 200

	for trial := 0; trial < ITERS; trial++ {
		server := &networking.Server{
			Hosts: []string{"./host-a", "./host-b", "*/otherhost", "host-c"},
		}

		var wg sync.WaitGroup
		wg.Add(N * 2)

		for i := 0; i < N; i++ {
			// Writer: writes server.Hosts[i] inside gateway.go (upstream file)
			go func() {
				defer wg.Done()
				sanitizeServerHostNamespace(server, "test-ns")
			}()
			// Concurrent writer racing on the same slice element
			go func() {
				defer wg.Done()
				sanitizeServerHostNamespace(server, "other-ns")
			}()
		}
		wg.Wait()
	}
}
