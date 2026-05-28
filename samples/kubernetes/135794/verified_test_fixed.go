package volumemanager

import (
	"errors"
	"sync"
	"testing"
)

// FIX version: mutex-protected append on errorsList eliminates the data race.
func TestRace_135794(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 50
	iters := 200

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				var errorsList []error
				var mu sync.Mutex
				var innerWg sync.WaitGroup
				innerWg.Add(2)
				go func() {
					defer innerWg.Done()
					mu.Lock()
					errorsList = append(errorsList, errors.New("err1"))
					mu.Unlock()
				}()
				go func() {
					defer innerWg.Done()
					mu.Lock()
					errorsList = append(errorsList, errors.New("err2"))
					mu.Unlock()
				}()
				innerWg.Wait()
				_ = errorsList
			}
		}()
	}

	wg.Wait()
}
