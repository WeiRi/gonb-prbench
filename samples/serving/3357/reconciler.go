package reconciler

// Base mirrors the BUG state of knative/serving pkg/reconciler/reconciler.go
// (PR #3357). NewBase spawns a logger broadcaster goroutine with no shutdown
// signal, outliving subtests. The goroutine continuously increments a shared
// counter (count) without synchronization — equivalent to the original bug
// where the logger goroutine reads t.common during/when the test teardown.
//
// Test goroutines also call Increment / GetCount on the same shared field,
// producing concurrent read/write races on count.

type Base struct {
	count  int          // BUG: unsynchronized shared field
	stopCh chan struct{}
}

// NewBase creates a Base and starts a background goroutine that continuously
// increments count without synchronization — modeling the leaked broadcaster.
func NewBase() *Base {
	b := &Base{stopCh: make(chan struct{})}
	go func() {
		for {
			select {
			case <-b.stopCh:
				return
			default:
				b.count++ // RACE write — goroutine outlives test scope
			}
		}
	}()
	return b
}

// Increment writes to the shared count field without synchronization.
func (b *Base) Increment() {
	b.count++ // RACE write
}

// GetCount reads the shared count field without synchronization.
func (b *Base) GetCount() int {
	return b.count // RACE read
}

// Stop signals the background goroutine to exit.
func (b *Base) Stop() {
	close(b.stopCh)
}
