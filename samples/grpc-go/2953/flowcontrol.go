package transport

// trInFlow BUG state: fields written without mutex (PR #2953 adds sync.Mutex)
type trInFlow struct {
	limit   uint32
	unacked uint32
}

func (f *trInFlow) newLimit(limit uint32) {
	f.limit = limit // BUG: write without lock → RACE
}

func (f *trInFlow) onData(n uint32) {
	f.unacked += n // BUG: write without lock → RACE
}
