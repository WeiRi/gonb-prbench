// Production stub for thanos pkg/pool/pool.go (PR #1340).
// Pre-PR: Put expression ^uint64(p.usedTotal-1) reads usedTotal non-atomically;
// Get does atomic.AddUint64.
package pool

import "sync/atomic"

type BytesPool struct {
	minSize    int
	maxSize    int
	growFactor int
	maxTotal   uint64
	usedTotal  uint64
}

func NewBytesPool(minSize, maxSize, growFactor int, maxTotal uint64) (*BytesPool, error) {
	return &BytesPool{minSize: minSize, maxSize: maxSize, growFactor: growFactor, maxTotal: maxTotal}, nil
}

// Get atomically increments usedTotal.
func (p *BytesPool) Get(sz int) (*[]byte, error) {
	atomic.AddUint64(&p.usedTotal, uint64(sz))
	b := make([]byte, sz)
	return &b, nil
}

// Put reads p.usedTotal NON-ATOMICALLY in the expression (BUG).
func (p *BytesPool) Put(b *[]byte) {
	sz := uint64(len(*b))
	// RACE: bare read of p.usedTotal-sz inside the AddUint64 argument expression
	atomic.AddUint64(&p.usedTotal, ^uint64(p.usedTotal-sz)+1)
}
