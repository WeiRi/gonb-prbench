package certprovider

import (
	"context"
	"sync"
	"testing"
)

type testProvider7128 struct{}
func (testProvider7128) KeyMaterial(context.Context) (*KeyMaterial, error) { return &KeyMaterial{}, nil }
func (testProvider7128) Close() {}

func TestRace_7128_InPlace(t *testing.T) {
	st := &store{providers: make(map[storeKey]*wrappedProvider)}
	sk := storeKey{name: "test"}
	p := &testProvider7128{}

	const N = 50
	const ITERS = 100
	for trial := 0; trial < ITERS; trial++ {
		wp := &wrappedProvider{Provider: p, storeKey: sk, store: st, refCount: 2}
		st.providers[sk] = wp
		var wg sync.WaitGroup
		wg.Add(N * 2)
		for i := 0; i < N; i++ {
			go func() {
				defer wg.Done()
				wp.Close()
			}()
			go func() {
				defer wg.Done()
				_, _ = wp.KeyMaterial(context.Background())
			}()
		}
		wg.Wait()
	}
}
