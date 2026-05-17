package pss

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: HandshakeController.releaseKey reads/writes ctl.symKeyIndex map without lock.
// FIX (PR #19162): releaseKey wraps with ctl.lock; the body is split off as releaseKeyNoLock.
// Same test code: in BUG concurrent releaseKey races on map; in FIX it's serialized.
func TestRace_19162_releaseKey(t *testing.T) {
	ctl := &HandshakeController{
		symKeyIndex: make(map[string]*handshakeKey),
	}
	tp := Topic{1}
	keys := make([]string, 20)
	for i := range keys {
		k := fmt.Sprintf("key-%d", i)
		keys[i] = k
		ctl.symKeyIndex[k] = &handshakeKey{symKeyID: &k}
	}

	var done int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			ctl.releaseKey(keys[i%len(keys)], &tp)
		}
		atomic.StoreInt32(&done, 1)
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && atomic.LoadInt32(&done) == 0; i++ {
			ctl.releaseKey(keys[(i+5)%len(keys)], &tp)
		}
	}()
	wg.Wait()
}
