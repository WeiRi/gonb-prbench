package rafthttp

import "sync"

// Local types.ID stub
type ID uint64

type Peer interface {
	stop()
}

type fakePeer struct{}

func (fakePeer) stop() {}

func newFakePeer() Peer { return fakePeer{} }

type Transport struct {
	mu    sync.Mutex
	peers map[ID]Peer
}

func (t *Transport) Pause() {
	for _, p := range t.peers {
		_ = p
	}
}

func (t *Transport) Resume() {
	for _, p := range t.peers {
		_ = p
	}
}
