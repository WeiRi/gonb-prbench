package logging

import (
	"sync"
	"sync/atomic"
	"testing"
)

// FIX: closedLock renamed to fileLock, purgeOldFiles writes oldestLogFileIdx under fileLock.
func TestRace_14538_oldestLogFileIdx(t *testing.T) {
	f := &FileRotator{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			f.fileLock.Lock()
			f.oldestLogFileIdx = i
			f.fileLock.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			f.fileLock.Lock()
			_ = f.oldestLogFileIdx
			f.fileLock.Unlock()
		}
	}()
	wg.Wait()
}
