package leaderelection

import (
	"sync"
	"testing"
	"time"

	rl "k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/utils/clock"
	clocktesting "k8s.io/utils/clock/testing"
)

func TestRace_136068(t *testing.T) {
	le := &LeaderElector{
		clock:              clocktesting.NewFakeClock(time.Now()),
		observedRecordLock: sync.Mutex{},
	}

	var wg sync.WaitGroup
	n := 50

	// Readers: read le.observedTime WITHOUT lock
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				_ = le.observedTime
			}
		}()
	}

	// Writers: call setObservedRecord which writes le.observedTime under lock
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rec := &rl.LeaderElectionRecord{
				HolderIdentity:       "test",
				LeaseDurationSeconds: 15,
			}
			for j := 0; j < 200; j++ {
				le.setObservedRecord(rec)
			}
		}()
	}

	wg.Wait()
}

var _ = clock.RealClock{}
