// Race-trigger test for grpc-go-1897; see README.md for usage.

package grpc

import (
	"sync"
	"testing"
	"time"
)

func TestRace_PR1897_MinConnectTimeout(t *testing.T) {
	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = resetTransport(stop)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			minConnectTimeout = time.Millisecond * 500
			time.Sleep(time.Microsecond)
			minConnectTimeout = time.Second * 20
		}
	}()

	time.Sleep(50 * time.Millisecond)
	close(stop)
	wg.Wait()
}
