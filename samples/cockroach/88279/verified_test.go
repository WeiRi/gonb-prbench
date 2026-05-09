package admission

import (
	"sync"
	"testing"
)

func Test88279Race(t *testing.T) {
	g := newElasticCPUGranter()
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			switch idx % 3 {
			case 0:
				_ = g.tryGet(1)
			case 1:
				g.tookWithoutPermission(1)
			case 2:
				g.setUtilizationLimit(0.5)
			}
		}(i)
	}
	wg.Wait()
}
