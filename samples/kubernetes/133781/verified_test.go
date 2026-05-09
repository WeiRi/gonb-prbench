package v1

import (
	"sync"
	"testing"
)

func TestRace_133781(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(2)
		// Writer: simulates encoding that modifies returned PriorityClasses
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				classes := SystemPriorityClasses()
				for _, c := range classes {
					c.Kind = "PriorityClass"
				}
			}
		}()
		// Reader: reads PriorityClasses fields (through IsKnownSystemPriorityClass)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				IsKnownSystemPriorityClass("system-node-critical", 2000001000, false)
			}
		}()
	}

	wg.Wait()
}
