// Race-trigger test for grpc-go-3090; see README.md for usage.

package grpc

import (
	"sync"
	"testing"
)

func TestRace_PR3090_ResolverWrapperBuild(t *testing.T) {
	const N = 200
	var wg sync.WaitGroup
	stops := make([]chan struct{}, N)
	for i := 0; i < N; i++ {
		stops[i] = make(chan struct{})
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			_ = NewCCResolverWrapper(stops[i])
		}()
	}
	for _, s := range stops {
		close(s)
	}
	wg.Wait()
}
