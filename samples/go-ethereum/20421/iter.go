// Minimal extraction of pre-fix go-ethereum p2p/enode/iter.go for PR #20421.
// Race: sliceIter.Node() reads it.nodes without holding it.mu;
// sliceIter.Close() writes it.nodes while holding it.mu.
// Production-code path: p2p/enode/iter.go
package enode

import "sync"

// Node mirrors the upstream identity stub.
type Node struct{ id int }

// sliceIter mirrors the pre-fix struct shape.
type sliceIter struct {
	mu    sync.Mutex
	nodes []*Node
	index int
	cycle bool
}

// Next is unchanged by the PR (already locked).
// Upstream path: p2p/enode/iter.go (pre-fix)
func (it *sliceIter) Next() bool {
	it.mu.Lock()
	defer it.mu.Unlock()
	if len(it.nodes) == 0 {
		return false
	}
	it.index++
	if it.index >= len(it.nodes) {
		if !it.cycle {
			it.nodes = nil
			return false
		}
		it.index = 0
	}
	return true
}

// Node — pre-fix version: NO mutex. PR #20421 adds it.mu.Lock/Unlock here.
// Upstream path: p2p/enode/iter.go (pre-fix line ~88-94).
func (it *sliceIter) Node() *Node {
	if len(it.nodes) == 0 {
		return nil
	}
	return it.nodes[it.index]
}

// Close already holds the lock pre-fix. Upstream path: p2p/enode/iter.go.
func (it *sliceIter) Close() {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.nodes = nil
}
