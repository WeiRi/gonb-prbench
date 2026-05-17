package registrytest

import (
	"context"
	"sync"
	"testing"

	"errors"
)

func TestRace_84487_WatchNodesRace(t *testing.T) {
	r := &NodeRegistry{}
	var wg sync.WaitGroup
	const N = 200
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			r.SetError(errors.New("e"))
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			_, _ = r.WatchNodes(context.Background(), nil)
		}
	}()
	wg.Wait()
}
