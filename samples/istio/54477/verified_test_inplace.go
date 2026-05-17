// Regression test for istio#54477
// PR: https://github.com/istio/istio/pull/54477
// Race: concurrent write to debugHandlers map via addDebugHandler
//        vs concurrent read of debugHandlers map via Debug.
// Both access sites are in debug.go (upstream production code).
// Uses real upstream types: DiscoveryServer, real upstream methods: addDebugHandler, Debug.
package xds

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRace_54477_InPlace(t *testing.T) {
	s := &DiscoveryServer{
		debugHandlers: make(map[string]string),
	}
	mux := http.NewServeMux()

	const N = 50
	const ITERS = 100

	for trial := 0; trial < ITERS; trial++ {
		var wg sync.WaitGroup
		wg.Add(N * 2)

		for i := 0; i < N; i++ {
			// Writer: calls addDebugHandler which writes s.debugHandlers[path] = help
			// This write happens inside debug.go (upstream file at line ~229)
			go func(id int) {
				defer wg.Done()
				s.addDebugHandler(mux, nil,
					fmt.Sprintf("/debug/race-%d-%d", trial, id),
					"help text",
					func(w http.ResponseWriter, r *http.Request) {},
				)
			}(i)

			// Reader: calls Debug which iterates over s.debugHandlers
			// This read happens inside debug.go (upstream file at line ~892)
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/debug", nil)
				s.Debug(w, r)
			}()
		}
		wg.Wait()
	}
}
