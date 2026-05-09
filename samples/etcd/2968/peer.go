package rafthttp

import "sync"

type msgAppReader struct{ id int }

type peer struct {
	mu           sync.Mutex
	msgAppReader *msgAppReader
	stopc        chan struct{}
}

// startPeer — BUG (pre-PR2968): returns p before the goroutine that initializes
// p.msgAppReader has run. External callers reading p.msgAppReader race with this
// init.
func startPeer() *peer {
	p := &peer{stopc: make(chan struct{})}
	go func() {
		// goroutine init writes msgAppReader (line 43)
		p.msgAppReader = &msgAppReader{id: 1}
		<-p.stopc
	}()
	return p
}

// stop — BUG: reads p.msgAppReader without sync (line 60).
func (p *peer) stop() *msgAppReader {
	close(p.stopc)
	return p.msgAppReader
}
