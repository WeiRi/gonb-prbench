// Race-trigger test for grpc-go-8541; see README.md for usage.

package xdsclient

import (
	"sync"
	"testing"
)

func TestRace_PR8541_LazyLrsClientInit(t *testing.T) {
	const N = 50
	for i := 0; i < N; i++ {
		c := NewClientImpl()
		var wg sync.WaitGroup
		const G = 10
		for j := 0; j < G; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = c.ReportLoad()
			}()
		}
		wg.Wait()
	}
}
