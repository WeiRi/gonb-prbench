package raft

type Message struct{ From, To uint64 }
type ConfChange struct{ ID uint64 }
type ConfState struct{ Voters []uint64 }
type Status struct{ Foo int }

type node struct {
	propc      chan Message
	recvc      chan Message
	confc      chan ConfChange
	confstatec chan ConfState
	done       chan struct{}
	tickc      chan struct{}
	status     chan chan Status
	stop       chan struct{}
}

// ClearPropc — BUG (pre-PR5165): writes n.propc without sync (line 21).
func (n *node) ClearPropc() {
	n.propc = nil // line 21 BUG
}

// ReadPropc — BUG: reads n.propc without sync (line 26).
func (n *node) ReadPropc() chan Message {
	return n.propc // line 26 BUG
}
