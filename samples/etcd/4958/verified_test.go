package v3rpc

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRace_4958(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 50
	iterations := 200

	// Writers: use atomic.StoreInt32 (as tests do)
	for g := 0; g < numGoroutines/2; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				atomic.StoreInt32(&ProgressReportIntervalMilliseconds, int32(1000+i))
			}
		}(g)
	}

	// Readers: plain read (as sendLoop does) - THIS IS THE RACE
	for g := 0; g < numGoroutines/2; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				// Plain read of the int32 - race with atomic.StoreInt32 above
				interval := time.Duration(ProgressReportIntervalMilliseconds) * time.Millisecond
				_ = interval
			}
		}(g)
	}

	wg.Wait()
}
