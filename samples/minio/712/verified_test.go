package disk

import (
	"sync"
	"testing"
)

// TestRaceGetFSInfo triggers data races on the buggy Disk struct
// which lacks a sync.Mutex. Concurrent GetFSInfo calls race on fsInfo map writes.
// Fix (PR #712): add sync.Mutex to Disk with Lock/Unlock in all methods.
func TestRaceGetFSInfo(t *testing.T) {
	numGoroutines := 80
	numIterations := 500

	// Construct Disk manually with fsInfo map already allocated
	disk := Disk{
		path:   "/tmp",
		fsInfo: make(map[string]string),
	}

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				// GetFSInfo writes to disk.fsInfo map WITHOUT lock
				// RACE: concurrent map writes from multiple goroutines
				_ = disk.GetFSInfo()
			}
		}()
	}
	wg.Wait()
}
