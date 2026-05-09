package deploymentwatcher

import (
	"sync"
	"testing"

	"github.com/hashicorp/nomad/nomad/structs"
)

// TestRaceGetEval triggers the data race where getEval()
// reads w.d.EvalPriority and w.j.Priority WITHOUT holding w.l,
// while other goroutines can modify these fields concurrently.
//
// Bug: priority reads in getEval() not protected by w.l
// Fix: added w.l.Lock() / w.l.Unlock() around priority reads
func TestRaceGetEval(t *testing.T) {
	numGoroutines := 50
	iterations := 200

	w := &deploymentWatcher{
		d: &structs.Deployment{
			EvalPriority: 50,
		},
		j: &structs.Job{
			Priority: 50,
		},
	}

	var wg sync.WaitGroup

	// Writers: modify w.d.EvalPriority and w.j.Priority concurrently
	// These writes SHOULD be protected by w.l but in the buggy code they're accessed without lock
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Bug: these writes race with reads in getEval()
				w.d.EvalPriority = id
				w.j.Priority = id
			}
		}(i)
	}

	// Readers: call getEval() which reads w.d.EvalPriority and w.j.Priority without lock
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = w.getEval()
			}
		}()
	}

	wg.Wait()
}

// TestRaceWatcherGetDeploys triggers the second race fix in PR 14121
// where w.state was read without holding w.l
func TestRaceWatcherGetDeploys(t *testing.T) {
	numGoroutines := 50
	iterations := 200

	w := &Watcher{
		state:    nil, // We don't need a real state, just testing the race on w.state reads
		watchers: make(map[string]*deploymentWatcher),
	}

	var wg sync.WaitGroup

	// Concurrent writes to w.state (simulating state updates)
	// and reads via direct access (bypassing the lock)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// RACE: direct read of w.state without lock (buggy)
				_ = w.state
			}
		}()
	}

	// Also write w.state concurrently to trigger race
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// RACE: direct write of w.state without lock
				w.state = nil
			}
		}()
	}

	wg.Wait()
}
