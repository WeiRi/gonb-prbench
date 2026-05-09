// Production stub for kubernetes-136068.
// Pre-PR: observedTime is read without lock (RACE) while setObservedRecord
// writes observedTime under observedRecordLock.
package leaderelection

import (
	"sync"
	"time"

	rl "k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/utils/clock"
)

type LeaderElector struct {
	clock              clock.Clock
	observedRecordLock sync.Mutex
	observedTime       time.Time
	observedRecord     rl.LeaderElectionRecord
}

// setObservedRecord writes observedTime under observedRecordLock.
func (le *LeaderElector) setObservedRecord(record *rl.LeaderElectionRecord) {
	le.observedRecordLock.Lock()
	defer le.observedRecordLock.Unlock()
	le.observedRecord = *record
	le.observedTime = le.clock.Now()
}
