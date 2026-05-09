package pacer

import (
	"sync"
	"testing"
	"time"
)

func TestRace_140_PacerFields(t *testing.T) {
	p := New()
	const N = 200
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			p.SetMinSleep(time.Duration(i+1) * time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			p.SetRetries(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			p.endCall(i%2 == 0)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			_ = p.Call()
		}
	}()
	wg.Wait()
}
