// Production stub for nats-server server/events.go (PR #1178).
// Pre-PR: debugSubscribers reads `nsubs` bare while sublist callbacks
// atomically AddInt32 -> race on int32.
package buggy

import (
	"sync/atomic"
)

type Server struct {
	results map[string]int32
}

// debugSubscribers mirrors the racy pre-PR path: a callback adds to nsubs
// concurrently with a deferred bare read.
func (s *Server) debugSubscribers(reply string) {
	var nsubs int32
	done := make(chan struct{})
	for k := 0; k < 4; k++ {
		go func() {
			atomic.AddInt32(&nsubs, 1)
			done <- struct{}{}
		}()
	}
	defer func() {
		// RACE: bare read of nsubs (line 1260 pre-PR) racing with atomic write.
		s.results[reply] = nsubs
	}()
	for k := 0; k < 4; k++ {
		<-done
	}
}
