// Race test for pq-1000 — conn.bad bool field race
package pq

import (
	"sync"
	"testing"
)

func TestRace_1000_ConnBadField(t *testing.T) {
	cn := newTestConn()
	const N = 200
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			markBad(cn)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			_ = isBad(cn)
		}
	}()
	wg.Wait()
}
