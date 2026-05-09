package group

import (
	"sync"
	"testing"
)

func TestRace_109969_GroupAdderSharedBackingArray(t *testing.T) {
	const N = 200
	// Caller owns a *shared* slice with extra capacity (e.g. from a cache).
	for n := 0; n < N; n++ {
		base := make([]string, 0, 16)
		base = append(base, "user1", "user2")
		// Two concurrent goroutines pass DIFFERENT Response structs but with
		// User.Groups all pointing to the SAME shared backing array.
		r1 := &Response{User: &DefaultInfo{Name: "a", Groups: base}}
		r2 := &Response{User: &DefaultInfo{Name: "b", Groups: base}}
		var wg sync.WaitGroup
		ag := &AuthenticatedGroupAdder{}
		wg.Add(2)
		go func() { defer wg.Done(); ag.AuthenticateRequest(r1) }()
		go func() { defer wg.Done(); ag.AuthenticateRequest(r2) }()
		wg.Wait()
	}
}
