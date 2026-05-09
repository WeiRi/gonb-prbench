// PR #19162 - swarm/pss/handshake.go - data race on
// HandshakeController.symKeyIndex map. Pre-fix: handler() and
// SendSym() read/write symKeyIndex map without ctl.lock; releaseKey
// also mutates it under different lock paths.
// PR fix: getSymKey() / registerSymKeyUse() / releaseKeyNoLock all
// take ctl.lock around symKeyIndex access.
// Production-code path: swarm/pss/handshake.go
package pss

import "sync"

type handshakeKey struct {
	count int
	limit int
}

type HandshakeController struct {
	lock         sync.Mutex
	symKeyIndex  map[string]*handshakeKey
}

func NewHandshakeController() *HandshakeController {
	return &HandshakeController{
		symKeyIndex: make(map[string]*handshakeKey),
	}
}

func (ctl *HandshakeController) Insert(id string) {
	ctl.lock.Lock()
	defer ctl.lock.Unlock()
	ctl.symKeyIndex[id] = &handshakeKey{count: 0, limit: 1000}
}

// Pre-fix handler: reads ctl.symKeyIndex[symkeyid] WITHOUT taking ctl.lock.
// Upstream: swarm/pss/handshake.go (pre-fix line ~298-312).
func (ctl *HandshakeController) Handler(symkeyid string) bool {
	if ctl.symKeyIndex[symkeyid] != nil {
		if ctl.symKeyIndex[symkeyid].count >= ctl.symKeyIndex[symkeyid].limit {
			return false
		}
		ctl.symKeyIndex[symkeyid].count++
	}
	return true
}

// Pre-fix releaseKey: mutates ctl.symKeyIndex (delete) under ctl.lock,
// but Handler reads without lock — write-vs-read race on map.
// Upstream: swarm/pss/handshake.go (pre-fix line ~215-232).
func (ctl *HandshakeController) ReleaseKey(symkeyid string) bool {
	ctl.lock.Lock()
	defer ctl.lock.Unlock()
	if ctl.symKeyIndex[symkeyid] == nil {
		return false
	}
	delete(ctl.symKeyIndex, symkeyid)
	return true
}
