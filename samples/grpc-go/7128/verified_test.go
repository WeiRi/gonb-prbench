// Race-trigger test for grpc-go-7128; see README.md for usage.

package main

import ("sync"; "testing")

type provider_7128 struct { closed bool }
func (p *provider_7128) Close() { p.closed = true }

func TestRace_7128(t *testing.T) {
    const N = 50; const ITERS = 100
    for trial := 0; trial < ITERS; trial++ {
        w := &struct{p *provider_7128}{p: &provider_7128{}}
        var wg sync.WaitGroup; wg.Add(N)
        for i := 0; i < N; i++ { go func() { defer wg.Done()
            w.p.Close()
        }() }
        wg.Wait()
    }
}