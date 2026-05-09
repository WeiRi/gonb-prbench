package cache

import (
	"strconv"
	"sync"
	"testing"
)

// TestRace_PR6947_SizeUnlocked reproduces the data race fixed by
// https://github.com/etcd-io/etcd/pull/6947 — cache.Size() reads c.lru.items
// without holding c.mu while Add() writes via c.mu.Lock().
func TestRace_PR6947_SizeUnlocked(t *testing.T) {
	c := NewCache(2048)
	const N = 8
	const ITERS = 2000
	var wg sync.WaitGroup
	wg.Add(N * 2)
	for i := 0; i < N; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				c.Add("k"+strconv.Itoa(id)+"-"+strconv.Itoa(j), j)
			}
		}(i)
		go func() {
			defer wg.Done()
			for j := 0; j < ITERS; j++ {
				_ = c.Size() // racy read of c.lru.items
			}
		}()
	}
	wg.Wait()
}
