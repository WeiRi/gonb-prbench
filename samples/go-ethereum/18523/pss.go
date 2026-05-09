// PR #18523 - swarm/pss/pss.go - data race on topicHandlerCaps map.
// Pre-fix: Pss.topicHandlerCaps map written in Register/deregister and
// read in handlePssMsg without any synchronisation. Issue ethersphere/swarm#1157
// shows WARNING DATA RACE: write at pss.go:360 (deregister, mapassign) and
// read at pss.go:889 (send, mapaccess1).
// PR fix: add topicHandlerCapsMu sync.RWMutex with getTopicHandlerCaps /
// setTopicHandlerCaps accessor methods.
// Production-code path: swarm/pss/pss.go (pre-fix line ~360, ~889).
package pss

type Topic int
type handlerCaps struct {
	raw  bool
	prox bool
}

type Pss struct {
	topicHandlerCaps map[Topic]*handlerCaps
}

func NewPss() *Pss {
	return &Pss{topicHandlerCaps: make(map[Topic]*handlerCaps)}
}

// Pre-fix Register: writes p.topicHandlerCaps without lock.
// Upstream: swarm/pss/pss.go (pre-fix line ~340, 360).
func (p *Pss) Register(topic Topic, raw bool) {
	if _, ok := p.topicHandlerCaps[topic]; !ok {
		p.topicHandlerCaps[topic] = &handlerCaps{}
	}
	if raw {
		p.topicHandlerCaps[topic].raw = true
	}
}

// Pre-fix HandlePssMsg: reads p.topicHandlerCaps without lock.
// Upstream: swarm/pss/pss.go (pre-fix line ~889).
func (p *Pss) HandlePssMsg(topic Topic) bool {
	if c, ok := p.topicHandlerCaps[topic]; ok {
		return c.raw
	}
	return false
}
