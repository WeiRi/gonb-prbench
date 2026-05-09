package main

import (
    "sync"
    "testing"
)

type store_2953 struct {
    m map[string]int
}

func new_2953() *store_2953 {
    return &store_2953{m: make(map[string]int)}
}

func (s *store_2953) set_2953(k string, v int) {
    s.m[k] = v
}

func (s *store_2953) get_2953(k string) int {
    return s.m[k]
}

func TestRace_2953(t *testing.T) {
    obj := new_2953()
    const N = 50
    const ITERS = 100
    var wg sync.WaitGroup
    wg.Add(N * 2)
    for i := 0; i < N; i++ {
        go func() {
            defer wg.Done()
            for j := 0; j < ITERS; j++ {
                obj.set_2953(string(rune(j)), j)
            }
        }()
    }
    for i := 0; i < N; i++ {
        go func() {
            defer wg.Done()
            for j := 0; j < ITERS; j++ {
                _ = obj.get_2953(string(rune(j)))
            }
        }()
    }
    wg.Wait()
}