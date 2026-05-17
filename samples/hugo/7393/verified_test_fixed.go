package hugolib

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX (PR #7393): h.getContentMaps() uses sync.Once to lazy-init h.content.
// The content field still exists, but production callers use the method.
func TestRace_7393_content_lazy_init(t *testing.T) {
	h := &HugoSites{numWorkers: 1}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = h.getContentMaps()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = h.getContentMaps()
		}
	}()
	wg.Wait()
}
