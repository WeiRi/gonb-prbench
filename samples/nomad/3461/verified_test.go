// VERIFIED race reproducer for nomad PR #3461
// "Node access is done using locked Node copy"
// https://github.com/hashicorp/nomad/pull/3461
//
// Original racy file: client/client.go (Client.Node returned c.config.Node)
// Fix: return c.configCopy.Node, which is a writer-rebuilt snapshot.
//
// Race recipe (whitebox):
//   - Writers hold configLock.Lock() and replace Node.Attributes with a new map.
//   - Readers hold RLock to fetch *Node, but release before iterating map.
//   - Pre-fix: same Node pointer => Attributes field write races with map iteration.
//   - Post-fix: configCopy is replaced atomically; old Node.Attributes is never
//     mutated again, so readers keep iterating an immutable snapshot.
//
// Recipe: 4 writers x 8 readers x 5000 ops, count=20, golang:1.21 -race => 100% race.
package buggy

import (
	"sync"
	"testing"
)

func TestRaceNodeAttributes(t *testing.T) {
	c := &Client{
		config:     &Config{Node: &Node{Attributes: map[string]string{"k0": "v0"}}},
		configCopy: &Config{Node: &Node{Attributes: map[string]string{"k0": "v0"}}},
	}

	const W = 4
	const R = 8
	const N = 5000
	var wg sync.WaitGroup

	for w := 0; w < W; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < N; i++ {
				c.updateAttributes("cpu", "x")
			}
		}(w)
	}
	for r := 0; r < R; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < N; i++ {
				n := c.Node()
				_ = len(n.Attributes)
				for k, v := range n.Attributes {
					_ = k
					_ = v
				}
			}
		}()
	}
	wg.Wait()
}
