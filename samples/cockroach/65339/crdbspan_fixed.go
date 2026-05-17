package tracing

import "sync"

// Reproduction of PR cockroachdb/cockroach#65339:
// "tracing: fix benign data race". BUG: crdbSpan.operation field
// is written by SetOperationName WITHOUT mu but read by getRecordingLocked
// holding mu.

type crdbSpan struct {
	mu        sync.Mutex
	operation string // BUG: not under mu in pre-fix code
	mu        struct {
		sync.Mutex
		duration int64
	}
}

// getRecordingLocked reads operation while caller holds mu.
func (s *crdbSpan) getRecordingLocked() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.operation // BUG line 21
}

// SetOperationName writes operation WITHOUT mu (BUG).
func (s *crdbSpan) SetOperationName(op string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.operation = op // BUG line 26
}

// GetRecording acquires mu and reads operation.
func (s *crdbSpan) GetRecording() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getRecordingLocked()
}
