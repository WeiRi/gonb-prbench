package main

import "sync"

type Shared struct {
	mu   sync.Mutex
	val  int64
	data map[string]int
}

func New() *Shared { return &Shared{data: make(map[string]int)} }

func (s *Shared) Write(v int64) { s.val = v; s.data["k"] = int(v) }
func (s *Shared) Read() int64 { return s.val }

// backendURLsBug models backendURLs in BUG state: stopped is a bare bool
type backendURLsBug struct {
	stopped bool // BUG: bare bool, race target
	wg      sync.WaitGroup
}

func (b *backendURLsBug) setBroken() {
	if b.stopped { // RACE READ on stopped
		return
	}
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		_ = b.stopped // RACE READ inside health-check goroutine
	}()
}

func (b *backendURLsBug) stopHealthChecks() {
	b.stopped = true // RACE WRITE on stopped
	b.wg.Wait()
}
