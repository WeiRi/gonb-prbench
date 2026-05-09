// Production stub for consul sdk/freeport/freeport.go (PR #7567).
// reset() modifies freePorts/pendingPorts/total without stopping the
// checkFreedPorts background goroutine -> race.
package freeport

import (
	"container/list"
	"errors"
	"time"
)

var (
	freePorts    *list.List
	pendingPorts *list.List
	total        int
	initDone     bool
)

func initialize() {
	freePorts = list.New()
	pendingPorts = list.New()
	total = 100
	for i := 0; i < total; i++ {
		freePorts.PushBack(20000 + i)
	}
	if !initDone {
		initDone = true
		go checkFreedPorts()
	}
}

func reset() {
	// racy: writes shared state without locking out checkFreedPorts goroutine
	freePorts = list.New()
	pendingPorts = list.New()
	total = 0
}

func Take(n int) ([]int, error) {
	if freePorts == nil {
		return nil, errors.New("not initialized")
	}
	out := make([]int, 0, n)
	for i := 0; i < n; i++ {
		e := freePorts.Front()
		if e == nil {
			return out, errors.New("no free ports")
		}
		freePorts.Remove(e)
		out = append(out, e.Value.(int))
	}
	return out, nil
}

func Return(ports []int) {
	if pendingPorts == nil {
		return
	}
	for _, p := range ports {
		pendingPorts.PushBack(p)
	}
}

func checkFreedPorts() {
	for {
		if pendingPorts == nil {
			return
		}
		// read pendingPorts/freePorts without lock — race vs reset()
		for e := pendingPorts.Front(); e != nil; e = e.Next() {
			freePorts.PushBack(e.Value)
		}
		time.Sleep(time.Microsecond)
	}
}
