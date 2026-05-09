package server

import (
	"sync"
	"testing"
)

// TestRaceMirrorSetupConsumerWg reproduces the data race in setupMirrorConsumer
// (server/stream.go) where mirror.wg.Wait() at line ~3304 races with concurrent
// mirror.wg.Add() at line ~3372 and mirror.wg.Done() in processMirrorMsgs.
//
// The bug: setupMirrorConsumer() launches a go func(){} closure that captures
// the 'mirror' variable. wg.Wait() at 3304 uses the captured mirror's wg.
// When processMirrorMsgs runs concurrently and calls wg.Done(), and another
// setupMirrorConsumer goroutine calls wg.Add(1) at 3372, all three operations
// (Wait, Add, Done) can race on the sync.WaitGroup's internal state because
// none of them hold mset.mu or any other lock.
//
// Fix: capture &mirror.wg into mirrorWg BEFORE the goroutine closure creation
// at line 3145, then use mirrorWg.Wait() instead of mirror.wg.Wait() at 3304.
// This pins the WaitGroup to the specific sourceInfo struct captured at entry.
func TestRaceMirrorSetupConsumerWg(t *testing.T) {
	// Create shared sourceInfo like mset.mirror
	si := &sourceInfo{}
	si.qch = make(chan struct{})

	numGoroutines := 100
	iterations := 300
	var allDone sync.WaitGroup
	ready := make(chan struct{})

	// Goroutines that simulate mirror.wg.Wait() at line 3304
	// These READ the WaitGroup's internal state without holding mset.mu.
	for i := 0; i < numGoroutines/2; i++ {
		allDone.Add(1)
		go func(id int) {
			defer allDone.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				si.wg.Wait()
			}
		}(i)
	}

	// Goroutines that simulate mirror.wg.Add(1) at line 3372 and
	// mirror.wg.Done() called by processMirrorMsgs goroutine.
	// These WRITE the WaitGroup's internal state without holding mset.mu.
	for i := 0; i < numGoroutines/2; i++ {
		allDone.Add(1)
		go func(id int) {
			defer allDone.Done()
			<-ready
			for j := 0; j < iterations; j++ {
				si.wg.Add(1)
				si.wg.Done()
			}
		}(i)
	}

	close(ready)
	allDone.Wait()
}
