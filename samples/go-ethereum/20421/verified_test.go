package enode

import (
	"sync"
	"testing"
)

// Race on sliceIter: Node() reads it.nodes without lock in BUG; Close() writes nodes=nil under lock.
func TestRace_20421_sliceIter_Node_vs_Close(t *testing.T) {
	nodes := make([]*Node, 32)
	for i := range nodes {
		nodes[i] = &Node{}
	}
	it := CycleNodes(nodes).(*sliceIter)
	// Drive Next once to set up state.
	it.Next()
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
		for i := 0; i < 500; i++ {
			it.Close()
		}
	}()
	wg.Wait()
}
