package consul1214

import "sync"

// ChildProcess models the shared child process reference
type ChildProcess struct {
	pid int
}

// LockCommand replicates the racy LockCommand from command/lock.go
type LockCommand struct {
	childLock sync.Mutex
	child     *ChildProcess
}

// startChild writes c.child under lock (safe)
func (c *LockCommand) startChild() {
	child := &ChildProcess{pid: 42}

	c.childLock.Lock()
	c.child = child
	c.childLock.Unlock()
}

// killChild reads c.child WITHOUT holding the lock (RACY - buggy version)
// In the buggy code, the lock was released before this function was called.
func (c *LockCommand) killChild() {
	// BUG: reading c.child without holding childLock
	// In the original code, killChild accessed c.child without the lock
	child := c.child

	if child != nil {
		_ = child.pid
	}
}
