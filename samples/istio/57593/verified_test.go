// Regression test for istio#57593
// PR: https://github.com/istio/istio/pull/57593
package main
import ("sync"; "testing")
type T_57593 struct { val int64 }
func (t *T_57593) write(v int64) { t.val = v }
func (t *T_57593) read() int64 { return t.val }
func TestRace_57593(t *testing.T) {
    obj := &T_57593{}
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