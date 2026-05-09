package garbagecollector

import (
	"fmt"
	"io"
	"sync"
)

// Stripped reproduction of pkg/controller/garbagecollector/graph.go pre-PR #134379.
// BUG: String() locks only dependentsLock; %#v reads beingDeleted/virtual fields
// that are written by markBeingDeleted/markObserved under different (or no) locks.

type OwnerReference struct {
	APIVersion, Kind, Name, UID string
}

type objectReference struct {
	OwnerReference
}

type node struct {
	identity        objectReference
	dependentsLock  sync.RWMutex
	dependents      map[*node]struct{}
	beingDeleted    bool   // written by markBeingDeleted, read by String via %#v (BUG)
	virtual         bool   // written by markObserved, read by String via %#v (BUG)
}

// String — BUG: locks ONLY dependentsLock; %#v reflects beingDeleted/virtual without sync.
func (n *node) String() string {
	n.dependentsLock.RLock()
	defer n.dependentsLock.RUnlock()
	return fmt.Sprintf("node:%#v", n)   // line 28 — reflective read of all fields
}

// Discard interface guard
var _ = io.Discard

func (n *node) markBeingDeleted() {
	n.beingDeleted = true               // line 36 — racing write
}

func (n *node) markObserved() {
	n.virtual = false                   // line 40 — racing write
}
