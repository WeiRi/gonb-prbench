// Regression test for kubernetes#107748
// PR: https://github.com/kubernetes/kubernetes/pull/107748
package main
import ("sync"; "testing")
type T_107748 struct { val int64 }
func (t *T_107748) write(v int64) { t.val = v }
func (t *T_107748) read() int64 { return t.val }
func TestRace_107748(t *testing.T) {
    obj := &T_107748{}
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