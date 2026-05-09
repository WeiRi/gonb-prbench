// Pre-fix topologycache.go from PR #118189 (endpointslice/topologycache).
// BUG: AddHints reads/writes hintsPopulatedByService and endpointsByService
// (lines 174 / 183 / 265) without t.lock; concurrent SetHints/HasPopulatedHints
// callers race on those map fields.
package topologycache

import "sync"

type ServiceSet struct {
	m map[string]struct{}
}

func newServiceSet() *ServiceSet {
	return &ServiceSet{m: map[string]struct{}{}}
}

func (s *ServiceSet) Has(k string) bool {
	_, ok := s.m[k]
	return ok
}

func (s *ServiceSet) Insert(k string) {
	s.m[k] = struct{}{}
}

type TopologyCache struct {
	lock                    sync.Mutex
	hintsPopulatedByService *ServiceSet
	endpointsByService      map[string]map[string]int // svc -> addrType -> count
}

func NewTopologyCache() *TopologyCache {
	return &TopologyCache{
		hintsPopulatedByService: newServiceSet(),
		endpointsByService:      map[string]map[string]int{},
	}
}

// AddHints — topologycache.go:174 in pre-fix: reads hintsPopulatedByService
// WITHOUT lock, then SetHints under lock => race window on map header.
func (t *TopologyCache) AddHints(serviceKey, addrType string) bool {
	hintsEnabled := t.hintsPopulatedByService.Has(serviceKey) // line 174 racy READ
	t.SetHints(serviceKey, addrType)
	return hintsEnabled
}

// SetHints — topologycache.go:183: writes under lock.
func (t *TopologyCache) SetHints(serviceKey, addrType string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.endpointsByService[serviceKey]; !ok {
		t.endpointsByService[serviceKey] = map[string]int{}
	}
	t.endpointsByService[serviceKey][addrType]++
	t.hintsPopulatedByService.Insert(serviceKey) // also writes the set
}

// HasPopulatedHints — topologycache.go:265 pre-fix: reads WITHOUT lock.
func (t *TopologyCache) HasPopulatedHints(serviceKey string) bool {
	return t.hintsPopulatedByService.Has(serviceKey) // racy READ of set.m
}
