package stats

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestRace_1317(t *testing.T) {
	ls := NewLeaderStats("leader1")

	var wg sync.WaitGroup
	numGoroutines := 50
	iterations := 200

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				// Simulate the pattern from send() in cluster_store.go:
				// direct access to ls.Followers map without lock
				to := strconv.FormatUint(uint64(i%10), 16)
				fs, ok := ls.Followers[to]
				if !ok {
					fs = &FollowerStats{}
					fs.Latency.Minimum = 1 << 63
					ls.Followers[to] = fs
				}

				// Call methods on FollowerStats concurrently
				if i%2 == 0 {
					fs.Succ(time.Millisecond)
				} else {
					fs.Fail()
				}
			}
		}(g)
	}
	wg.Wait()
}
