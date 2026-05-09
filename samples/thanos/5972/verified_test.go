// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package query

import (
	"sync"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// TestRace_5972 triggers the data race on e.endpoints map between
// Close() (writer, holds Lock) and Update()/getTimedOutRefs() (readers, no RLock).
//
// BUG: EndpointSet.Close() acquires e.endpointsMtx.Lock() and writes to
// e.endpoints. EndpointSet.Update() reads e.endpoints at line ~396 WITHOUT
// e.endpointsMtx.RLock(). Similarly, getTimedOutRefs() reads e.endpoints
// at line ~455 without RLock.
//
// FIX: Adds e.endpointsMtx.RLock()/RUnlock() around the read paths.
//
// Level 1 WARNING DATA RACE artifact in PR body with full stack traces.
func TestRace_5972(t *testing.T) {
	reg := prometheus.NewRegistry()

	es := NewEndpointSet(
		func() time.Time { return time.Now() },
		log.NewNopLogger(),
		reg,
		func() []*GRPCEndpointSpec { return nil },
		[]grpc.DialOption{grpc.WithInsecure()},
		5*time.Minute,
		5*time.Second,
	)

	// Pre-populate endpoints map so readers have data.
	dummyRef := &endpointRef{
		addr:    "127.0.0.1:10901",
		created: time.Now(),
		status:  &EndpointStatus{LastCheck: time.Now()},
		logger:  log.NewNopLogger(),
	}
	es.endpointsMtx.Lock()
	es.endpoints["127.0.0.1:10901"] = dummyRef
	es.endpointsMtx.Unlock()

	numGoroutines := 50
	iterations := 200

	var wg sync.WaitGroup

	// Reader goroutines: call getTimedOutRefs() which reads e.endpoints
	// WITHOUT RLock in the BUG state.
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = es.getTimedOutRefs()
			}
		}()
	}

	// Writer goroutine: simulates Close() by writing to e.endpoints under Lock.
	// This races with getTimedOutRefs() which reads without RLock.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < iterations; j++ {
			es.endpointsMtx.Lock()
			es.endpoints = map[string]*endpointRef{
				"127.0.0.1:10901": dummyRef,
			}
			es.endpointsMtx.Unlock()
		}
	}()

	// Also exercise the direct read path from Update() (line ~396):
	// for addr, er := range e.endpoints { ... }
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < iterations; j++ {
			// Directly simulate the unprotected read in Update().
			for addr, er := range es.endpoints {
				_, _ = addr, er
			}
		}
	}()

	wg.Wait()
}
