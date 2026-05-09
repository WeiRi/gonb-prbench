package channelz

import "sync"

type Channel interface {
	ChannelzMetric() string
}

type realChannel struct{ v int32 }

func (r *realChannel) ChannelzMetric() string { return "real" }

type dummyChannel struct{}

func (dummyChannel) ChannelzMetric() string { return "dummy" }

type cn struct {
	c Channel
}

type ChannelMap struct {
	mu sync.RWMutex
	m  map[int64]*cn
}

func NewChannelMap() *ChannelMap {
	return &ChannelMap{m: map[int64]*cn{}}
}

func (m *ChannelMap) Add(id int64, c Channel) {
	m.mu.Lock()
	m.m[id] = &cn{c: c}
	m.mu.Unlock()
}

// GetChannel — BUG (pre-PR2411): reads cn.c AFTER unlock, racing with Delete write.
func (m *ChannelMap) GetChannel(id int64) string {
	m.mu.RLock()
	c := m.m[id]
	if c == nil {
		m.mu.RUnlock()
		return ""
	}
	m.mu.RUnlock()
	return c.c.ChannelzMetric() // line 69 BUG: read after unlock
}

func (m *ChannelMap) DeleteSelfFromMap(id int64) {
	m.mu.Lock()
	if c, ok := m.m[id]; ok {
		c.c = dummyChannel{} // line 84 write under lock
	}
	m.mu.Unlock()
}
