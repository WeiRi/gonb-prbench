// Regression test for grpc-go#8605
// PR: https://github.com/grpc/grpc-go/pull/8605
// RACE_TARGET: ADS flow control — pending flag + channel replaced with sync.Cond
// FVM: FVM-4 — complex sync restructuring with Cond
package xdsclient
import ("sync"; "testing")
type flowCtrl_8605 struct { mu sync.Mutex; pending bool }
func (fc *flowCtrl_8605) wait() { fc.mu.Lock(); fc.mu.Unlock() }
func (fc *flowCtrl_8605) setPending(v bool) { fc.pending = v }
func TestRace_8605_InPlace(t *testing.T) {
    fc := &flowCtrl_8605{}
    const N = 50; const ITERS = 100
    var wg sync.WaitGroup; wg.Add(N * 2)
    for i := 0; i < N; i++ { go func() { defer wg.Done()
        for j := 0; j < ITERS; j++ { fc.wait() }
    }() }
    for i := 0; i < N; i++ { go func() { defer wg.Done()
        for j := 0; j < ITERS; j++ { fc.setPending(j%2==0) }
    }() }
    wg.Wait()
}