package v1

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Race: BUG SystemPriorityClasses() returns the underlying slice
// (shallow), so callers share pointers to elements. Concurrent
// callers that read/write PriorityClass.Description race on the
// PriorityClass.Description field.
// FIX returns a deep-copy of each element.
func TestRace_kubernetes_133781_priorityclasses_shared(t *testing.T) {
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 50000 && atomic.LoadInt32(&done) == 0; j++ {
			pcs := SystemPriorityClasses()
			pcs[0].Description = "mut-a"
		}
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 50000 && atomic.LoadInt32(&done) == 0; j++ {
			pcs := SystemPriorityClasses()
			_ = pcs[0].Description
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
