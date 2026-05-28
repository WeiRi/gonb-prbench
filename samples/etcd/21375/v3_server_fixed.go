package etcdserver

import "errors"

var ErrLeaderChanged = errors.New("etcdserver: leader changed")

type readState struct {
	Index uint64
}

type EtcdServer struct{}

// FIX: after readStateC fires, re-check leaderChangedNotifier in inner select.
func (s *EtcdServer) requestCurrentIndex(leaderChangedNotifier <-chan struct{}, readStateC <-chan readState) (uint64, error) {
	select {
	case rs := <-readStateC:
		select {
		case <-leaderChangedNotifier:
			return 0, ErrLeaderChanged
		default:
			return rs.Index, nil
		}
	case <-leaderChangedNotifier:
		return 0, ErrLeaderChanged
	}
}
