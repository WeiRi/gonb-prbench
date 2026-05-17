package server

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_nats_5182_mset_lseq(t *testing.T) {
	mset := &stream{}
	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			mset.mu.Lock()
			mset.lseq = uint64(j)
			mset.mu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 100000 && atomic.LoadInt32(&done) == 0; j++ {
			_ = mset.lastSeq()
		}
	}()
	wg.Wait()
}
