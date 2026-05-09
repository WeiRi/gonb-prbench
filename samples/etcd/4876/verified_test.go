package clientv3

import (
	"sync"
	"testing"
)

// TestRace_PR4876_SwitchRemoteUnlocked reproduces the data race fixed by
// https://github.com/etcd-io/etcd/pull/4876 — switchRemote reads kv.conn
// before locking, racing with Do() readers and concurrent switchRemote() writers.
func TestRace_PR4876_SwitchRemoteUnlocked(t *testing.T) {
	k := newKV()
	const N = 8
	const ITERS = 2000
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = k.Do()
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				k.switchRemote()
			}
		}()
	}
	wg.Wait()
}
