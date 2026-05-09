// Race-trigger test for grpc-go-8202; see README.md for usage.

package delegatingresolver

import (
	"sync"
	"testing"
)

func TestRace_PR8202_DelegatingResolverChildMu(t *testing.T) {
	const N = 200
	for i := 0; i < N; i++ {
		r := New()
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			r.updateProxyResolverState()
			r.ResolveNow()
		}()
		go func() {
			defer wg.Done()
			r.Close()
		}()
		wg.Wait()
	}
}
