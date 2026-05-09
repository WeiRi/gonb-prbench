package filters

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// TestRace_107452_HeaderMutation reproduces the data race from PR #107452.
//
// Race scenario (BUG state):
//   - Handler goroutine: calls Header().Set() which modifies the underlying
//     writer's header map (returned by tw.w.Header())
//   - Timeout goroutine: calls WriteHeaderTimeout() which writes to the
//     underlying writer directly (tw.w.WriteHeader), reading the header map
//   - Both goroutines access the SAME http.Header map concurrently -> data race
//
// In the FIX state, Header() returns a private handlerHeaders copy (cloned
// at construction), so the handler goroutine writes to handlerHeaders while
// the timeout goroutine reads from w.Header() — different maps, no race.
func TestRace_107452_HeaderMutation(t *testing.T) {
	const N = 1000
	const G = 12

	var wg sync.WaitGroup
	for g := 0; g < G; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < N; i++ {
				rec := httptest.NewRecorder()
				tw := &baseTimeoutWriter{w: rec}

				var inner sync.WaitGroup
				inner.Add(2)

				// Handler goroutine: write headers via Header().Set()
				// In BUG state, this modifies w.Header() (the underlying map)
				// In FIX state, this modifies handlerHeaders (private copy)
				go func() {
					defer inner.Done()
					for j := 0; j < 30; j++ {
						tw.Header().Set("X-Custom-Header", "value")
					}
				}()

				// Timeout goroutine: writes directly to underlying writer
				// This reads w.Header() to write HTTP status line + headers
				go func() {
					defer inner.Done()
					tw.WriteHeaderTimeout(http.StatusGatewayTimeout)
				}()

				inner.Wait()
			}
		}()
	}
	wg.Wait()
}
