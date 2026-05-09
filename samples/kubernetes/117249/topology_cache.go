// Production stub for kubernetes-117249.
// Pre-PR: getAllocations does the nil check on cpuRatiosByZone BEFORE acquiring
// lock; SetNodes writes cpuRatiosByZone under lock. RACE.
package topologycache

import (
	"sync"

	v1 "k8s.io/api/core/v1"
	discovery "k8s.io/api/discovery/v1"
)

type SliceInfo struct {
	ServiceKey  string
	AddressType discovery.AddressType
	ToCreate    []*discovery.EndpointSlice
}

type TopologyCache struct {
	lock             sync.Mutex
	cpuRatiosByZone  map[string]float64
}

func NewTopologyCache() *TopologyCache { return &TopologyCache{} }

// SetNodes writes cpuRatiosByZone under lock.
func (c *TopologyCache) SetNodes(nodes []*v1.Node) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cpuRatiosByZone = map[string]float64{}
	for _, n := range nodes {
		zone := n.Labels[v1.LabelTopologyZone]
		c.cpuRatiosByZone[zone] = float64(len(n.Status.Allocatable))
	}
}

// AddHints reads cpuRatiosByZone WITHOUT lock — RACE.
func (c *TopologyCache) AddHints(_ *SliceInfo) {
	c.getAllocations()
}

func (c *TopologyCache) getAllocations() map[string]float64 {
	if c.cpuRatiosByZone == nil { // RACE: read of map header without lock
		return nil
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	out := make(map[string]float64, len(c.cpuRatiosByZone))
	for k, v := range c.cpuRatiosByZone {
		out[k] = v
	}
	return out
}
