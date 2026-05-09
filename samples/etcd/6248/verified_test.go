// Regression test for etcd#6248
// PR: https://github.com/etcd/etcd/pull/6248
// FVM: SharedVar
package main
import ("sync"; "testing")
type T_6248 struct { val int64 }
func (t *T_6248) write(v int64) { t.val = v }
func (t *T_6248) read() int64 { return t.val }
func TestRace_6248(t *testing.T) {
    obj := &T_6248{}
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