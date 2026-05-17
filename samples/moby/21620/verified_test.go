//go:build linux
// +build linux

package libcontainerd

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

// BUG: client.lock(id) reads containerMutexes map AFTER Unlock; concurrent
// goroutines write the map → race on map.
func TestRace_moby_21620_libcontainerd_lock(t *testing.T) {
	c := &client{
		clientCommon: clientCommon{
			containers:       map[string]*container{},
			containerMutexes: map[string]*sync.Mutex{},
		},
	}
	var done int32
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			id := "id" + strconv.Itoa(gid%3)
			for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
				c.lock(id)
				c.unlock(id)
			}
			atomic.StoreInt32(&done, 1)
		}(i)
	}
	wg.Wait()
}
