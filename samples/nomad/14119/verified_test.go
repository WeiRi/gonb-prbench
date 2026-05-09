package client

import (
	"sync"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/nomad/structs"
)

// TestRaceHeartbeatLastOk triggers the data race where watch()
// writes h.lastOk WITHOUT holding the lock (heartbeatstop.go:73),
// while setLastOk() writes it under Lock (heartbeatstop.go:126-130).
//
// Bug: h.lastOk = time.Now() in watch() (heartbeatstop.go:73)
// Fix: replaced with h.setLastOk(time.Now()) 
func TestRaceHeartbeatLastOk(t *testing.T) {
	numInstances := 30
	iterations := 200

	logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	var wg sync.WaitGroup

	for i := 0; i < numInstances; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			shutdownCh := make(chan struct{})
			h := newHeartbeatStop(
				func(id string) (AllocRunner, error) { return nil, nil },
				10*time.Second,
				logger,
				shutdownCh,
			)
			h.allocHookCh = make(chan *structs.Allocation, 1)
			h.allocInterval = map[string]time.Duration{}

			// Start watch() goroutine — it does h.lastOk = time.Now()
			// WITHOUT lock at heartbeatstop.go:73 (PRODUCTION CODE RACE)
			var innerWg sync.WaitGroup
			innerWg.Add(1)
			go func() {
				defer innerWg.Done()
				h.watch()
			}()

			// Call setLastOk concurrently — this does Lock + write + Unlock
			// at heartbeatstop.go:126-130 (PROTECTED)
			for j := 0; j < iterations; j++ {
				h.setLastOk(time.Now())
			}

			// Shutdown watch() 
			close(shutdownCh)
			innerWg.Wait()
		}()
	}

	wg.Wait()
}
