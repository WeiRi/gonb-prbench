// Production stub for etcd server/etcdserver/v3_server.go (PR #21375).
// Models EtcdServer.requestCurrentIndex: select picks readStateC vs
// leaderChangedNotifier non-deterministically. BUG: when both are ready
// at the same time, picking readStateC returns success even though leader
// has changed -> stale read.
package etcdserver

import "errors"

var ErrLeaderChanged = errors.New("etcdserver: leader changed")

type readState struct {
	Index uint64
}

type EtcdServer struct{}

// requestCurrentIndex: BUG. Single select doesn't re-check leaderChangedNotifier
// after readStateC fires.
func (s *EtcdServer) requestCurrentIndex(leaderChangedNotifier <-chan struct{}, readStateC <-chan readState) (uint64, error) {
	select {
	case rs := <-readStateC:
		return rs.Index, nil
	case <-leaderChangedNotifier:
		return 0, ErrLeaderChanged
	}
}
