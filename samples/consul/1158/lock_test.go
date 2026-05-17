package main

import ("sync"; "testing")

func TestRace(t *testing.T) {
	for i := 0; i < 50; i++ {
		l := New()
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); l.Unlock() }()
		go func() { defer wg.Done(); _ = l.IsHeld() }()
		wg.Wait()
	}
}
