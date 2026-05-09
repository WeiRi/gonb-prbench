package fixture

import (
	"sync"
	"testing"
)

func TestRace_103487_FixtureWatcher(t *testing.T) {
	const N = 200
	var wg sync.WaitGroup
	for n := 0; n < N; n++ {
		obj := &Object{Data: map[string]string{"key": "v"}}
		tr := NewTracker()
		tr.add("a", obj)

		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				_ = obj.Data["key"] // racy READ (other goroutine writes via Modify/Add/Delete)
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				tr.add("a", obj)
			}
		}()
	}
	wg.Wait()
}
