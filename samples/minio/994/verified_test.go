// Race test for minio-994 — contentdb.Init() concurrent calls race on extDB
// BUG: Init() does extDB = make(...) and then conditionally loadDB; no mutex
// FIX: mutex guards Init, plus double-check of isInitialized
package contentdb

import (
	"sync"
	"testing"
)

func TestRace_994_InitConcurrent(t *testing.T) {
	// Reset state to simulate cold start
	isInitialized = false
	extDB = nil
	const N = 30
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = Init()
		}()
	}
	wg.Wait()
}
