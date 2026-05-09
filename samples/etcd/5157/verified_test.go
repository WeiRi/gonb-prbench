package stats

import (
	"sync"
	"testing"
	"time"
)

func TestRace_5157(t *testing.T) {
	const numGoroutines = 50
	const iterations = 200

	for g := 0; g < numGoroutines; g++ {
		var wg sync.WaitGroup
		ss := &ServerStats{
			ID:    "node1",
			State: StateFollower,
		}
		ss.LeaderInfo.StartTime = time.Now()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				ss.SetState(StateLeader, "node1")
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_, _, _ = ss.GetState()
			}
		}()

		wg.Wait()
	}
}
