// Production stub for nats-server server/stream.go (PR #7395).
// Pre-PR: setupMirrorConsumer launches a goroutine that closes over mirror,
// using mirror.wg.Wait() at line 3304; concurrent calls Add/Done on the same
// wg without mset.mu -> wg internal state race.
package server

import "sync"

type sourceInfo struct {
	wg  sync.WaitGroup
	qch chan struct{}
}
