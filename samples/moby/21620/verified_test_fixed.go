//go:build linux
// +build linux

package libcontainerd

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/docker/docker/pkg/locker"
)

// FIX: clientCommon uses *locker.Locker (per-id locking with internal sync).
// Concurrent lock/unlock through the locker are safe.
func TestRace_moby_21620_libcontainerd_lock(t *testing.T) {
	c := &client{
		clientCommon: clientCommon{
			containers: map[string]*container{},
			locker:     locker.New(),
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
