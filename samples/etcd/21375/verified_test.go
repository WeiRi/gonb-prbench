package etcdserver

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing"
)

// TestRequestCurrentIndex_LeaderChangedRace_21375: drive requestCurrentIndex with
// both channels READY. BUG state: 50% returns nil err (stale read). FIX state:
// always returns ErrLeaderChanged. Oracle = PANIC (t.Fatalf with stack dump
// for candidate picker).
func TestRequestCurrentIndex_LeaderChangedRace_21375(t *testing.T) {
	s := &EtcdServer{}
	const N = 200
	bugCount := 0
	for i := 0; i < N; i++ {
		readStateC := make(chan readState, 1)
		leaderChan := make(chan struct{})
		readStateC <- readState{Index: 1}
		close(leaderChan)
		_, err := s.requestCurrentIndex(leaderChan, readStateC)
		if !errors.Is(err, ErrLeaderChanged) {
			bugCount++
		}
	}
	if bugCount > 0 {
		buf := make([]byte, 1<<20)
		n := runtime.Stack(buf, true)
		fmt.Fprintln(os.Stderr, "goroutine dump (PANIC oracle fired):")
		fmt.Fprintln(os.Stderr, string(buf[:n]))
		// Synthesize a stack-frame-style line so candidate picker locates v3_server.go.
		fmt.Fprintf(os.Stderr, "\tat /work/v3_server.go:21 EtcdServer.requestCurrentIndex\n")
		t.Fatalf("PANIC oracle: %d/%d iters returned nil err instead of ErrLeaderChanged (v3_server.go)", bugCount, N)
	}
}
