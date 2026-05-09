// Production stub for swarmkit agent/node.go (PR #1046).
// Pre-PR: runManager launches an inner goroutine that reads `ready` via closure.
// After outer select exits, runManager writes `ready = nil` -> race with inner read.
package agent

import (
	"context"
	"time"
)

// RunManagerLoopBUG mirrors the BUG-state runManager() inner-goroutine race.
func RunManagerLoopBUG(ctx context.Context, done chan struct{}) {
	ready := make(chan struct{})
	go func() {
		time.Sleep(time.Microsecond) // ensure inner runs while outer writes nil
		select {
		case <-ready: // RACE: read of captured 'ready' slot
		case <-ctx.Done():
		}
	}()
	select {
	case <-ctx.Done():
	case <-done:
	}
	ready = nil // RACE: write to captured slot
	_ = ready
}
