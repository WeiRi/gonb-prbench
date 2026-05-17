package layer

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_moby_21677_layer_get_refcount(t *testing.T) {
	ls := &layerStore{
		layerMap: map[ChainID]*roLayer{},
		mounts:   map[string]*mountedLayer{},
	}
	cid := ChainID("sha256:test")
	ls.layerMap[cid] = &roLayer{
		chainID:        cid,
		referenceCount: 1,
		references:     map[Layer]struct{}{},
		layerStore:     ls,
	}

	var done int32
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 2000 && atomic.LoadInt32(&done) == 0; j++ {
				_, _ = ls.Get(cid)
			}
			atomic.StoreInt32(&done, 1)
		}()
	}
	wg.Wait()
}
