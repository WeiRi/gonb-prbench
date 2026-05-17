package query

import (
	"sync"
	"testing"
	"time"
)

func TestRace_5972_EndpointsMapIter(t *testing.T) {
	es := &EndpointSet{
		endpoints:                make(map[string]*endpointRef),
		now:                      time.Now,
		unhealthyEndpointTimeout: time.Hour,
	}
	es.endpoints["a"] = &endpointRef{addr: "a", created: time.Now()}
	var wg sync.WaitGroup
	const N = 200
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			es.endpointsMtx.Lock()
			es.endpoints["b"] = &endpointRef{addr: "b", created: time.Now()}
			delete(es.endpoints, "b")
			es.endpointsMtx.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			_ = es.getTimedOutRefs()
		}
	}()
	wg.Wait()
}
