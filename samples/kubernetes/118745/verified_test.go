// Race-trigger test for kubernetes-118745; see README.md for usage.

package controller

import (
	"io"
	"sync"
	"testing"
)

func TestExpectationsLoggingRace(t *testing.T) {
	exp := &ControlleeExpectations{add: 1, del: 0, key: "ns/foo"}
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			exp.Add(1, 0) // atomic.AddInt64
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			logExpectations(io.Discard, "exp:", exp) // BUG: reflective %#v
		}
	}()

	wg.Wait()
}
