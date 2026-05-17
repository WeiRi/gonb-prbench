package channelz

import (
	"sync"
	"sync/atomic"
	"testing"
)

// fakeChannel implements Channel for race test.
type fakeChannel struct{}

func (fakeChannel) ChannelzMetric() *ChannelInternalMetric { return &ChannelInternalMetric{} }

// Race: BUG channelMap.GetChannel reads `cn.c.ChannelzMetric()` AFTER RUnlock;
// concurrent goroutine (via channel.deleteSelfFromMap) sets `cn.c = &dummyChannel{}`
// under Lock. The cn.c field read after RUnlock races with the field write.
// FIX captures `chanCopy := cn.c` BEFORE RUnlock.
func TestRace_grpc_go_2411_channelmap_chanCopy(t *testing.T) {
	cm := &channelMap{
		topLevelChannels: map[int64]struct{}{},
		servers:          map[int64]*server{},
		channels:         map[int64]*channel{},
		subChannels:      map[int64]*subChannel{},
		listenSockets:    map[int64]*listenSocket{},
		normalSockets:    map[int64]*normalSocket{},
	}
	const cid = int64(1)
	cm.channels[cid] = &channel{
		id:          cid,
		c:           fakeChannel{},
		nestedChans: map[int64]string{},
		subChans:    map[int64]string{},
		cm:          cm,
		trace:       &channelTrace{},
	}

	var done int32
	var wg sync.WaitGroup
	// reader: many GetChannel calls
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
				_ = cm.GetChannel(cid)
			}
		}()
	}
	// writer: flip cn.c between fakeChannel and dummyChannel
	wg.Add(1)
	go func() {
		defer wg.Done()
		cn := cm.channels[cid]
		for j := 0; j < 50000 && atomic.LoadInt32(&done) == 0; j++ {
			cm.mu.Lock()
			if j%2 == 0 {
				cn.c = &dummyChannel{}
			} else {
				cn.c = fakeChannel{}
			}
			cm.mu.Unlock()
		}
		atomic.StoreInt32(&done, 1)
	}()
	wg.Wait()
}
