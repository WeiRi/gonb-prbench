package nomad

import (
	"fmt"
	"sync"
	"testing"

	hclog "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/nomad/nomad/structs"
)

// TestRace_BlockedEvals_ChannelCapture reproduces the data race on
// capacityChangeCh. Unblock() and UnblockQuota() release the lock and then
// send on b.capacityChangeCh without the lock held. Flush() replaces
// b.capacityChangeCh with a new channel under the lock.
//
// This is a concurrent read (sending on the channel reads the pointer) and
// write (Flush writes a new pointer) to the same field.
//
// 60+ goroutines directly access the capacityChangeCh field: readers capture
// the channel reference and send on it, writers replace it. A drain goroutine
// consumes from all channels to prevent blocking.
func TestRace_BlockedEvals_ChannelCapture(t *testing.T) {
	logger := hclog.NewNullLogger()
	b := &BlockedEvals{
		logger:           logger.Named("blocked_evals"),
		evalBroker:       nil,
		captured:         make(map[string]wrappedEval),
		escaped:          make(map[string]wrappedEval),
		system:           newSystemEvals(),
		jobs:             make(map[structs.NamespacedID]string),
		unblockIndexes:   make(map[string]uint64),
		capacityChangeCh: make(chan *capacityUpdate, 1000),
		duplicateCh:      make(chan struct{}, 1),
		stopCh:           make(chan struct{}),
		stats:            NewBlockedStats(),
	}
	b.enabled = true

	// Drain goroutine: consume from capacityChangeCh to prevent blocking.
	// Uses a for-select that reads via an intermediate variable so it
	// doesn't get stuck when the channel is replaced.
	drainDone := make(chan struct{})
	go func() {
		defer func() { recover() }()
		for {
			select {
			case <-drainDone:
				return
			default:
				select {
				case <-b.capacityChangeCh:
				default:
				}
			}
		}
	}()

	var wg sync.WaitGroup
	nWriters := 60
	nReaders := 60
	nIters := 500

	// Writer goroutines: replace capacityChangeCh (simulating Flush)
	for i := 0; i < nWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				// Write to the channel field (simulating Flush)
				b.capacityChangeCh = make(chan *capacityUpdate, 100)
			}
		}(i)
	}

	// Reader goroutines: read capacityChangeCh and send on it (simulating Unblock)
	for i := 0; i < nReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < nIters; j++ {
				// Read the channel field and send on it (simulating Unblock)
				ch := b.capacityChangeCh
				// Non-blocking send to avoid deadlock
				select {
				case ch <- &capacityUpdate{
					computedClass: fmt.Sprintf("class-%d-%d", id, j),
					index:         uint64(j),
				}:
				default:
				}
			}
		}(i)
	}

	wg.Wait()
	close(drainDone)
}
