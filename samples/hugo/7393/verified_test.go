package hugolib

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: HugoSites.content is set/read without sync.Once → race on lazy init.
// PR #7393 adds contentInit sync.Once and h.getContentMaps() method.
func TestRace_7393_content_lazy_init(t *testing.T) {
	h := &HugoSites{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			h.content = nil
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5000 && atomic.LoadInt32(&done) == 0; i++ {
			_ = h.content
		}
	}()
	wg.Wait()
}
