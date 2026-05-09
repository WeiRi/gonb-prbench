package integration

import (
	"sync"
	"testing"

	"ase/etcd-6197/clientv3"
)

func TestRace_PR6197_ToGRPC(t *testing.T) {
	// Add a client to the map so toGRPC returns early without calling ActiveConnection
	earlyClient := &clientv3.Client{}
	proxies[earlyClient] = grpcAPI{}

	const N = 50
	const ITERS = 200
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				// toGRPC reads proxies map (in cluster_proxy.go) - the buggy code has no lock
				toGRPC(earlyClient)
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				// Write to global proxies map without lock (simulates newClientV3 in buggy code)
				proxies[&clientv3.Client{}] = grpcAPI{}
			}
		}()
	}
	wg.Wait()
}
