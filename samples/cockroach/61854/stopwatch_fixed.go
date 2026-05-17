package timeutil

import (
	"sync"
	"time"
)

// Stopwatch is a simplified reproducer of cockroach
// pkg/util/timeutil/stopwatch.go (BUG state) where Start/Stop/Elapsed
// concurrently read/write startTime and elapsed without synchronization.
type Stopwatch struct {
	mu        sync.Mutex
	startTime time.Time     // line 52 area: write in Start
	elapsed   time.Duration // line 61: write in Stop
	running   bool          // line 67: read/write in Stop
}

func NewStopwatch() *Stopwatch { return &Stopwatch{} }

func (s *Stopwatch) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.startTime = time.Now() // race write
	s.running = true
}

func (s *Stopwatch) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running { // race read
		return
	}
	s.elapsed += time.Since(s.startTime) // race read+write
	s.running = false                    // race write
}

func (s *Stopwatch) Elapsed() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running { // race read
		return s.elapsed + time.Since(s.startTime) // race read
	}
	return s.elapsed
}
