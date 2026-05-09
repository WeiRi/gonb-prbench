package memory

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestRaceExpireObjects triggers a data race in the buggy expireObjects()
// where len(memory.objectMetadata) is read WITHOUT holding the lock (line 504),
// while concurrent CreateObject() calls write to objectMetadata under lock (line 306).
// Fix (commit 3a1386165): move memory.lock.Lock() before the len() check.
func TestRaceExpireObjects(t *testing.T) {
	numGoroutines := 60
	numIterations := 100

	_, _, driver := Start(1024*1024*1024, time.Second)

	err := driver.CreateBucket("bucket", "private")
	if err != nil {
		t.Logf("CreateBucket error: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				key := fmt.Sprintf("rk-%d-%d", id, j)
				data := strings.NewReader(fmt.Sprintf("v-%d-%d", id, j))
				// CreateObject writes to objectMetadata under lock
				// RACE: expireObjects goroutine reads len(objectMetadata) without lock
				driver.CreateObject("bucket", key, "application/octet-stream", "", data)
			}
		}(i)
	}
	wg.Wait()
}
