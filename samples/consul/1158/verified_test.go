// consul-1158: lock.go: fix race condition
// Original race: cmd.Start() was NOT under childLock, but c.child assignment was.
// killChild could run after cmd.Start() but before c.child assignment,
// resulting in an un-killed child process. Fix: moved childLock.Lock() BEFORE cmd.Start().
// Original diff file: command/lock.go
// Original frame hits: command/lock.go:271 (child assignment under lock, Start outside)

package consul1158

import (
	"sync"
	"testing"
)

// LockCommand replicates the racy LockCommand from command/lock.go
type LockCommand struct {
	childLock sync.Mutex
	child     *ChildProcess
	started   bool // Simulates cmd.Start() - not protected by childLock (RACY)
}

type ChildProcess struct {
	pid int
}

// startChild: cmd.Start() modifies .started (racy), then lock protects .child assignment
func (c *LockCommand) startChild() {
	// Simulates cmd.Start() - NOT under lock in buggy version
	c.started = true

	// childLock acquired AFTER Start (BUG)
	c.childLock.Lock()
	c.child = &ChildProcess{pid: 42}
	c.childLock.Unlock()
}

// killChild reads c.child under lock, but c.started is read without lock
func (c *LockCommand) killChild() {
	// BUG: c.started may be true (process running) but c.child may be nil
	// because startChild hasn't reached the Lock() section yet
	c.childLock.Lock()
	child := c.child
	c.childLock.Unlock()

	if child != nil {
		_ = child.pid
	}
}

// TestRace reproduces the race between startChild and killChild.
// 60 goroutines concurrently start and kill - c.started races with c.child.
func TestRace(t *testing.T) {
	iterations := 300

	for i := 0; i < iterations; i++ {
		lc := &LockCommand{}

		var wg sync.WaitGroup

		for g := 0; g < 30; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				lc.startChild()
			}()
		}

		for g := 0; g < 30; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				lc.killChild()
			}()
		}

		wg.Wait()
	}
}
