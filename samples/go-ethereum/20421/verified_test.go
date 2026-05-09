package enode

import (
	"sync"
	"testing"
)

// TestRace_20421_sliceIter_Node_vs_Close: races Node() (read nodes) vs
// Close() (writes nodes=nil under lock). PR #20421 fixes by adding mu.Lock
// to Node(). Production-code path: p2p/enode/iter.go.
func TestRace_20421_sliceIter_Node_vs_Close(t *testing.T) {
	nodes := make([]*Node, 32)
	for i := range nodes {
		nodes[i] = &Node{id: i}
	}
	it := &sliceIter{nodes: nodes, cycle: true}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			_ = it.Node()
		}
	}()
	go func() {
		defer wg.Done()
		// Drive Close to write the racy field.
		for i := 0; i < 500; i++ {
			it.Close()
		}
	}()
	wg.Wait()
}
