package client

import (
	"sync"
	"testing"
	"time"

	"ase/etcd-20656/model"
)

func TestRace_20656(t *testing.T) {
	const numGoroutines = 50
	const numIterations = 200

	for iter := 0; iter < numIterations; iter++ {
		c := &RecordingClient{baseTime: time.Now()}
		c.InitWatch(model.WatchRequest{Key: "test-key"})

		var wg sync.WaitGroup
		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				c.AppendRespUnsafe(0, model.WatchResponse{
					Revision: int64(id),
					Time:     time.Duration(id),
				})
			}(g)
		}
		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = c.ReadRespsUnsafe(0)
			}()
		}
		wg.Wait()
	}
}
