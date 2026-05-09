// Regression test for istio#54477
// PR: https://github.com/istio/istio/pull/54477
package main
import ("sync"; "testing")
type T_54477 struct { val int64 }
func (t *T_54477) write(v int64) { t.val = v }
func (t *T_54477) read() int64 { return t.val }
func TestRace_54477(t *testing.T) {
    obj := &T_54477{}
    const N = 50; const ITERS = 100
    var wg sync.WaitGroup; wg.Add(N * 2)
    for i := 0; i < N; i++ { go func() { defer wg.Done()
        for j := 0; j < ITERS; j++ { obj.write(int64(j)) }
    }() }
    for i := 0; i < N; i++ { go func() { defer wg.Done()
        for j := 0; j < ITERS; j++ { _ = obj.read() }
    }() }
    wg.Wait()
}