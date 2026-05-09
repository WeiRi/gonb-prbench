// Race-trigger test for grpc-go-3373; see README.md for usage.

package grpctest

import (
	"regexp"
	"sync"
	"testing"
)

func TestRace_PR3373_TLoggerMap(t *testing.T) {
	g := NewTLogger()
	g.AddExpect(regexp.MustCompile("foo"), 1000000)
	g.AddExpect(regexp.MustCompile("bar"), 1000000)

	var wg sync.WaitGroup
	const G = 6
	const N = 5000
	for i := 0; i < G; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				_ = g.expected("foobar test message")
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 50; j++ {
			g.Update()
			g.AddExpect(regexp.MustCompile("foo"), 1000)
			g.AddExpect(regexp.MustCompile("bar"), 1000)
		}
	}()
	wg.Wait()
}
