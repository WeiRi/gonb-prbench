package leaderelection

import (
	"sync"
	"testing"
	"time"

	"k8s.io/utils/clock"
	rl "k8s.io/client-go/tools/leaderelection/resourcelock"
)

func TestRace_136068_ObservedTimeRace(t *testing.T) {
	le := &LeaderElector{clock: clock.RealClock{}}
	var wg sync.WaitGroup
	const N = 200
	wg.Add(2)
	// Goroutine A: setObservedRecord writes observedTime under lock
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			le.observedRecordLock.Lock()
			le.observedTime = time.Now()
			le.observedRecord = rl.LeaderElectionRecord{LeaseDurationSeconds: 10}
			le.observedRecordLock.Unlock()
		}
	}()
	// Goroutine B: call isLeaseValid which reads observedTime (BUG: outside lock)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			_ = le.isLeaseValid(time.Now())
		}
	}()
	wg.Wait()
}
