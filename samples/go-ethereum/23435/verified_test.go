package miner

import (
	"testing"
	"time"
)

// TestRace_23435_worker_close_vs_mainLoop_current: close() reads
// w.current.state from external goroutine while mainLoop() writes w.current
// in its own goroutine — racy pointer + field access.
func TestRace_23435_worker_close_vs_mainLoop_current(t *testing.T) {
	w := newWorker()
	go w.MainLoop()
	time.Sleep(5 * time.Millisecond)
	w.Close()
	time.Sleep(5 * time.Millisecond)
}
