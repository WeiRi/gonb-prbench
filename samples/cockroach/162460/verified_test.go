package sql

import (
	"sync"
	"testing"
)

func Test162460Race(t *testing.T) {
	p := &Planner{}
	dsp := &DistSQLPlanner{}
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				dsp.CreatePhysPlan(p)
			} else {
				_ = CheckRunner(p)
			}
		}(i)
	}
	wg.Wait()
}
