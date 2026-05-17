package main

import "sync"

type Lock struct {
	mu    sync.Mutex
	held  bool
	owner string
}

// BUG: Lock reads held without mutex
func (l *Lock) IsHeld() bool { return l.held }

// BUG: Unlock writes held without mutex  
func (l *Lock) Unlock() { l.held = false; l.owner = "" }

func New() *Lock { return &Lock{held: true, owner: "test"} }
