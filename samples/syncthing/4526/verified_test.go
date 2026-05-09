package connections

import (
	"sync"
	"testing"
)

func TestRace_4526(t *testing.T) {
	tgts := make([]dialTarget, 20)
	for i := range tgts {
		tgts[i] = dialTarget{priority: i, uri: "tcp://x"}
	}
	var wg sync.WaitGroup
	for k := 0; k < 50; k++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			DialParallel(DeviceID{}, tgts)
		}()
	}
	wg.Wait()
}
