package volumemanager

import (
	"errors"
	"sync"
	"testing"
)

func TestRace_135794(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				// Bug: errors slice is shared across goroutines without
				// synchronization. Concurrent append causes data race.
				// This mirrors the buggy pattern in WaitForAllPodsUnmount
				// where multiple goroutines do errors = append(errors, err).
				var errorsList []error
				var innerWg sync.WaitGroup
				innerWg.Add(2)
				go func() {
					defer innerWg.Done()
					errorsList = append(errorsList, errors.New("err1"))
				}()
				go func() {
					defer innerWg.Done()
					errorsList = append(errorsList, errors.New("err2"))
				}()
				innerWg.Wait()
				_ = errorsList
			}
		}()
	}

	wg.Wait()
}
