package memory

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRaceExpireObjects(t *testing.T) {
	r := NewIntelligent(0, time.Hour)
	for i := 0; i < 100; i++ {
		r.Set(fmt.Sprintf("seed-%d", i), []byte("v"))
	}
	var wg sync.WaitGroup
	numGoroutines := 50
	iterations := 200
	wg.Add(numGoroutines * 2)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				r.Set(fmt.Sprintf("k-%d-%d", id, j), []byte("v"))
			}
		}(i)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				r.ExpireObjects(time.Microsecond)
			}
		}()
	}
	wg.Wait()
}
