package snapshot

import (
	"sync"
	"testing"
)

// TestRace_24685_diffLayer_Parent_vs_flatten: Parent() reads dl.parent
// without lock while flatten() writes under Lock — fires under -race.
func TestRace_24685_diffLayer_Parent_vs_flatten(t *testing.T) {
	dl := newDiffLayer(nil)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			dl.flatten(&diffLayer{})
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = dl.Parent()
		}
	}()
	wg.Wait()
}
