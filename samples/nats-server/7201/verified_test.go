// VERIFIED race reproducer for nats-server PR #7201
// "[FIXED] Hold consumer lock when reading o.cfg.PauseUntil"
// https://github.com/nats-io/nats-server/pull/7201
//
// Original racy file: server/jetstream_api.go (jsConsumerCreateRequest +
// jetstream_cluster.go processConsumerLeaderChange)
// Pre-fix: o.cfg.PauseUntil read without o.mu.RLock(); fix wraps with RLock/RUnlock.
//
// 4 writers x 8 readers x 5000 ops, golang:1.21 -race => race in 1st iter.
package buggy

import (
	"sync"
	"testing"
	"time"
)

func TestRacePauseUntilRead(t *testing.T) {
	now := time.Now().Add(time.Hour)
	s := &consumer{cfg: ConsumerConfig{PauseUntil: &now}}

	const W = 4
	const R = 8
	const N = 5000
	var wg sync.WaitGroup
	for i := 0; i < W; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < N; j++ {
				t1 := time.Now().Add(time.Duration(j) * time.Millisecond)
				s.updatePauseUntil(&t1)
			}
		}(i)
	}
	for i := 0; i < R; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := &ConsumerConfig{}
			for j := 0; j < N; j++ {
				_ = s.jsConsumerCreateRequest(req)
			}
		}()
	}
	wg.Wait()
}
