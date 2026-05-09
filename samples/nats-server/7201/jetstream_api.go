// Production stub for nats-server server/jetstream_api.go (PR #7201).
// Pre-PR: o.cfg.PauseUntil read without o.mu.RLock.
package buggy

import (
	"sync"
	"time"
)

type ConsumerConfig struct {
	PauseUntil *time.Time
}

type consumer struct {
	mu  sync.RWMutex
	cfg ConsumerConfig
}

// updatePauseUntil writes o.cfg.PauseUntil under lock.
func (o *consumer) updatePauseUntil(t *time.Time) {
	o.mu.Lock()
	o.cfg.PauseUntil = t
	o.mu.Unlock()
}

// jsConsumerCreateRequest reads o.cfg.PauseUntil WITHOUT RLock (pre-PR bug).
func (o *consumer) jsConsumerCreateRequest(req *ConsumerConfig) error {
	if o.cfg.PauseUntil != nil { // RACE: bare read
		_ = *o.cfg.PauseUntil
	}
	_ = req
	return nil
}
