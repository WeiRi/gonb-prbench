// Race-trigger test for serving-3068; see README.md for usage.

package pocserving3068

import (
	"errors"
	"sync"
)

// BUG version of knative/serving pkg/pool/pool.go (pre-PR-3068):
// Wait() closes workCh; concurrent Go() then sends on closed channel → panic.

type Interface interface {
	Go(func() error)
	Wait() error
}

type impl struct {
	wg     sync.WaitGroup
	workCh chan func() error
	errCh  chan error

	once sync.Once
	bail error
}

func NewWithCapacity(workers, capacity int) Interface {
	i := &impl{
		workCh: make(chan func() error, capacity),
		errCh:  make(chan error, capacity),
	}
	for w := 0; w < workers; w++ {
		go func() {
			for fn := range i.workCh {
				if err := fn(); err != nil {
					i.errCh <- err
				}
				i.wg.Done()
			}
		}()
	}
	return i
}

// Go is the BUG version: no doneCh guard. send on closed workCh panics.
func (i *impl) Go(w func() error) {
	i.wg.Add(1)
	i.workCh <- w // BUG: panics if workCh is closed by Wait()
}

func (i *impl) Wait() error {
	i.once.Do(func() {
		close(i.workCh) // BUG: closes before guaranteeing no more Go() calls
		go func() {
			i.wg.Wait()
			close(i.errCh)
		}()
		for err := range i.errCh {
			if err != nil && i.bail == nil {
				i.bail = err
			}
		}
	})
	return i.bail
}

var _ = errors.New
