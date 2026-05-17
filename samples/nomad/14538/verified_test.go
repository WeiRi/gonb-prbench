package logging

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: FileRotator.purgeOldFiles writes f.oldestLogFileIdx UNLOCKED while
// nextFile reads it under f.closedLock → race on oldestLogFileIdx.
// PR #14538 renames closedLock→fileLock and wraps the purgeOldFiles write.
func TestRace_14538_oldestLogFileIdx(t *testing.T) {
	f := &FileRotator{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			f.oldestLogFileIdx = i
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			f.closedLock.Lock()
			_ = f.oldestLogFileIdx
			f.closedLock.Unlock()
		}
	}()
	wg.Wait()
}
