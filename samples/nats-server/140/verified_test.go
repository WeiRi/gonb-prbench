package server

import (
	"sync/atomic"
	"testing"
)

func TestRaceGlobalTraceFlagReadWrite(t *testing.T) {
	numReaders := 80
	numWriters := 10
	iterations := 500
	done := make(chan struct{})
	ready := make(chan struct{})

	// Readers call client.traceOp which non-atomically reads global trace
	// in client.go:210 — the buggy production code path
	for i := 0; i < numReaders; i++ {
		go func(id int) {
			c := &client{}
			<-ready
			for j := 0; j < iterations; j++ {
				c.traceOp("test %s", "op", []byte("arg"))
			}
			done <- struct{}{}
		}(i)
	}

	// Writers toggle the global trace variable atomically
	// (same pattern as log.go that sets trace/debug)
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			<-ready
			for j := 0; j < iterations; j++ {
				atomic.StoreInt32(&trace, 1)
				atomic.StoreInt32(&trace, 0)
			}
			done <- struct{}{}
		}(i)
	}

	close(ready)

	for i := 0; i < numReaders+numWriters; i++ {
		<-done
	}
}
