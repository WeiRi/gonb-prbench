// Race-trigger test for etcd-6324; see README.md for usage.

package grpcproxy

import (
	"sync"
	"testing"
)

func TestRace_etcd_6324(t *testing.T) {
	const ITERS = 500
	for trial := 0; trial < 30; trial++ {
		sws := &serverWatchStream{singles: make(map[int64]*watcherSingle)}
		for i := int64(0); i < 3; i++ {
			sws.addDedicatedWatcher(i)
		}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := int64(100); i < 100+ITERS; i++ {
				sws.addDedicatedWatcher(i)
			}
		}()
		go func() {
			defer wg.Done()
			sws.close()
		}()
		wg.Wait()
	}
}
