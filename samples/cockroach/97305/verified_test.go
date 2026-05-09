package insights

import (
	"sync"
	"testing"
)

func Test97305Race(t *testing.T) {
	d := newAnomalyDetector()
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = d.GetPercentileValues(1)
		}()
	}
	wg.Wait()
}
