// Race-trigger test for grpc-go-3763; see README.md for usage.

package main

import ("sync"; "testing")

func TestRace_3763(t *testing.T) {
    var rpcSucceeded bool
    const N = 50; const ITERS = 100
    var wg sync.WaitGroup; wg.Add(N * 2)
    for i := 0; i < N; i++ { go func() { defer wg.Done()
        for j := 0; j < ITERS; j++ { rpcSucceeded = true }
    }() }
    for i := 0; i < N; i++ { go func() { defer wg.Done()
        for j := 0; j < ITERS; j++ { _ = rpcSucceeded }
    }() }
    wg.Wait()
}