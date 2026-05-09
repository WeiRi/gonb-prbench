// Copyright 2026 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package timeutil

import (
	"sync"
	"testing"
)

// TestRace_61854 exercises the data race in StopWatch where Elapsed() reads
// the elapsed field concurrently with Stop() writing to it, without any
// synchronization. The bug was fixed by wrapping all fields in a mutex.
func TestRace_61854(t *testing.T) {
	const (
		numGoroutines = 60
		numIterations = 300
	)

	for iter := 0; iter < numIterations; iter++ {
		w := NewStopwatch()

		var wg sync.WaitGroup

		// Goroutines that start/stop the watch concurrently.
		// Stop() writes to elapsed without synchronization (BUG).
		for g := 0; g < numGoroutines/2; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				w.Start()
				w.Stop()
			}()
		}

		// Goroutines that read elapsed concurrently.
		// Elapsed() reads elapsed without synchronization (BUG).
		for g := 0; g < numGoroutines/2; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = w.Elapsed()
			}()
		}

		wg.Wait()
	}
}
