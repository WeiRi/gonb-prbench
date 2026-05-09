// Race-trigger test for grpc-go-2411; see README.md for usage.

package channelz

import (
	"sync"
	"testing"
)

func TestRace_PR2411_ChannelzGetChannel(t *testing.T) {
	m := NewChannelMap()
	const N = 200
	for i := int64(0); i < N; i++ {
		m.Add(i, &realChannel{v: int32(i)})
	}

	var wg sync.WaitGroup
	const G = 6
	const ITERS = 5000

	for g := 0; g < G; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < ITERS; i++ {
				_ = m.GetChannel(int64(i % N))
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < ITERS; i++ {
			m.DeleteSelfFromMap(int64(i % N))
			m.Add(int64(i%N), &realChannel{v: int32(i)})
		}
	}()
	wg.Wait()
}
