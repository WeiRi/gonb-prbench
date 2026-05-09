package clientv3

import (
	"sync"
	"testing"
)

func TestRace_20192(t *testing.T) {
	c := &Client{
		lg:   NewNopLogger(),
		lgMu: &sync.RWMutex{},
	}

	const N = 20
	const ITERS = 10000

	var ready sync.WaitGroup
	ready.Add(N * 2)
	var start sync.WaitGroup
	start.Add(1)
	var done sync.WaitGroup
	done.Add(N * 2)

	for i := 0; i < N; i++ {
		go func() {
			ready.Done()
			start.Wait()
			defer done.Done()
			for j := 0; j < ITERS; j++ {
				c.SetLogger(NewNopLogger())
			}
		}()
		go func() {
			ready.Done()
			start.Wait()
			defer done.Done()
			for j := 0; j < ITERS; j++ {
				_ = c.GetLogger()
			}
		}()
	}

	ready.Wait()
	start.Done()
	done.Wait()
}
