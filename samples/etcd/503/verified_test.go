package lock

import (
	"sync"
	"testing"
)

func TestRace_503(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 50
	iterations := 200

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				// Reproduce the exact channel pattern from buggy watch():
				// stopWatchChan can be closed while goroutine tries to send
				stopWatchChan := make(chan bool)
				closeChan := make(chan bool)

				barrier := make(chan struct{})

				go func() {
					close(barrier) // signal we're about to enter select
					select {
					case <-closeChan:
						// RACE: sending on channel that may be closed
						stopWatchChan <- true
					case <-stopWatchChan:
					}
				}()

				<-barrier        // wait for goroutine to be ready
				close(closeChan) // trigger the closeChan case
				// RACE: close vs send
				close(stopWatchChan)
			}
		}()
	}

	wg.Wait()
}
