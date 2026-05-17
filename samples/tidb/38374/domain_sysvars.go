package main

import "sync"

type Shared struct {
	mu   sync.Mutex
	val  int64
	data map[string]int
}

func New() *Shared { return &Shared{data: make(map[string]int)} }

// BUG: write without lock
func (s *Shared) Write(v int64) { s.val = v; s.data["k"] = int(v) }

// BUG: read without lock  
func (s *Shared) Read() int64 { return s.val }
