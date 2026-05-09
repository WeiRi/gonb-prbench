package lease

import (
	"sync"
	"testing"
)

func Test144073Race(t *testing.T) {
	t1 := &descriptorState{}
	d := &descriptorVersionState{}
	d.refcount.Store(0)
	d.mu.lease = &storedLease{id: 1}
	t1.mu.active = []*descriptorVersionState{d}

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				_ = t1.removeInactiveVersions()
			} else {
				d.SetLease(&storedLease{id: idx})
			}
		}(i)
	}
	wg.Wait()
}
