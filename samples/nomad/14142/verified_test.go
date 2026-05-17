package nomad

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: Unblock reads b.capacityChangeCh AFTER releasing b.l, while Flush reassigns
// b.capacityChangeCh under b.l. Race on the field.
func TestRace_14142_capacityChangeCh(t *testing.T) {
	b := &BlockedEvals{
		enabled:          true,
		stats:            NewBlockedStats(),
		captured:         map[string]wrappedEval{},
		escaped:          map[string]wrappedEval{},
		unblockIndexes:   make(map[string]uint64),
		capacityChangeCh: make(chan *capacityUpdate, 1024),
		duplicateCh:      make(chan struct{}, 1),
		stopCh:           make(chan struct{}),
		system:           newSystemEvals(),
	}

	stopDrain := make(chan struct{})
	go func() {
		// Drain Unblock sends without reading b.capacityChangeCh (avoid concurrent
		// access to the field; the racy field access we want is in Unblock itself).
		for {
			select {
			case <-stopDrain:
				return
			default:
				// snapshot copy of the chan pointer inside b.l for safe read:
				b.l.RLock()
				ch := b.capacityChangeCh
				b.l.RUnlock()
				select {
				case <-ch:
				case <-stopDrain:
					return
				}
			}
		}
	}()
	defer close(stopDrain)

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			b.Unblock("class", uint64(i))
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 200 && atomic.LoadInt32(&done) == 0; i++ {
			b.Flush()
		}
	}()
	wg.Wait()
}
