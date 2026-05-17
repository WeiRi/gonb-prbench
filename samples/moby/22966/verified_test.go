package memory

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRace_moby_22966_register_vs_watch(t *testing.T) {
	d := &Discovery{}
	if err := d.Initialize("", 10*time.Millisecond, 0, nil); err != nil {
		t.Fatal(err)
	}

	stopCh := make(chan struct{})
	ch, errCh := d.Watch(stopCh)

	var done atomic.Bool
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for !done.Load() {
			select {
			case <-ch:
			case <-errCh:
			case <-time.After(5 * time.Millisecond):
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 500 && !done.Load(); i++ {
			_ = d.Register("addr")
		}
		done.Store(true)
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 500 && !done.Load(); i++ {
			_ = d.Register("addr2")
		}
	}()

	wg.Wait()
	close(stopCh)
}
