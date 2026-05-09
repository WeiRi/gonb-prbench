package vsphere

import "sync"

// Stripped reproduction of vsphere.go::VSphere.DisksAreAttached pre-PR #90348.
// BUG: outer "for dc, nodes := range dcNodes" + "go func() { ... attached, _ := nodeManager.checkAttached(nodes, dc) ... }()"
// captures loop vars dc and nodes; range loop overwrites them on next iteration.

type checkResult struct {
	dc   string
	mu   sync.Mutex
	hits map[string]bool
}

// ProcessDisksAreAttached mirrors the BUG-state goroutine launch.
func ProcessDisksAreAttached(dcNodes map[string][]string) map[string]bool {
	out := map[string]bool{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for dc, nodes := range dcNodes { // line 19 (range vars overwritten)
		wg.Add(1)
		go func() { // line 21 (captures dc & nodes by ref)
			defer wg.Done()
			// reads of captured loop vars while range loop overwrites them
			tag := dc + "/"           // line 24 (read of loop var dc)
			for _, n := range nodes { // line 25 (read of loop var nodes)
				mu.Lock()
				out[tag+n] = true
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return out
}
