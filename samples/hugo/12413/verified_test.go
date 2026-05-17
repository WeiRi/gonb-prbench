// Race test for hugo-12413 — hugocontext.Wrap returns shared bufferpool bytes
// BUG: Wrap returns buf.Bytes() then defer PutBuffer(buf) returns buf to pool.
//      Caller holds []byte aliasing pool memory; next Wrap reuses buf → race.
// FIX: Wrap returns buf.String() — immutable independent copy.
package hugocontext

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRace_hugo_12413_BufferAlias(t *testing.T) {
	const N = 100
	var wg sync.WaitGroup
	var counter int64
	for round := 0; round < 5; round++ {
		results := make([]interface{}, N)
		for i := 0; i < N; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = Wrap([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), uint64(idx))
			}(i)
		}
		wg.Wait()
		// Now race: reading prior results vs new Wraps (which reuse pool buffers)
		for i := 0; i < N; i++ {
			wg.Add(2)
			go func(r interface{}) {
				defer wg.Done()
				// Read prior result
				switch v := r.(type) {
				case []byte:
					if len(v) > 0 {
						atomic.AddInt64(&counter, int64(v[0]))
					}
				case string:
					if len(v) > 0 {
						atomic.AddInt64(&counter, int64(v[0]))
					}
				}
			}(results[i])
			go func(idx int) {
				defer wg.Done()
				_ = Wrap([]byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), uint64(idx+10000))
			}(i)
		}
		wg.Wait()
	}
	_ = counter
}
