package main

import (
    "sync"
    "testing"
)

type store_7070 struct {
    m map[string]int
}

func new_7070() *store_7070 {
    return &store_7070{m: make(map[string]int)}
}

func (s *store_7070) set_7070(k string, v int) {
    s.m[k] = v
}

func (s *store_7070) get_7070(k string) int {
    return s.m[k]
}

func TestRace_7070(t *testing.T) {
    obj := new_7070()
    const N = 50
    const ITERS = 100
    var wg sync.WaitGroup
    wg.Add(N * 2)
    for i := 0; i < N; i++ {
        go func() {
            defer wg.Done()
            for j := 0; j < ITERS; j++ {
                obj.set_7070(string(rune(j)), j)
            }
        }()
    }
    for i := 0; i < N; i++ {
        go func() {
            defer wg.Done()
            for j := 0; j < ITERS; j++ {
                _ = obj.get_7070(string(rune(j)))
            }
        }()
    }
    wg.Wait()
}