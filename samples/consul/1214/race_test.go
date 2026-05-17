package main

import ("sync"; "testing")

func TestRace(t *testing.T) {
	for i := 0; i < 50; i++ {
		s := New()
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); for j := 0; j < 200; j++ { s.Write(int64(j)) } }()
		go func() { defer wg.Done(); for j := 0; j < 200; j++ { _ = s.Read() } }()
		wg.Wait()
	}
}
