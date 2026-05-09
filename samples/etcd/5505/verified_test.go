package v3rpc

import (
	"sync"
	"testing"
	"time"
)

// TestRace_PR5505_CtrlChanClose reproduces the race fixed by
// https://github.com/etcd-io/etcd/pull/5505 — close(ctrlStream) in close()
// races with recvLoop()'s send. We swallow any send-on-closed panic so the
// test goroutine is allowed to finish; the race detector still flags the
// happens-before violation in close+send.
func TestRace_PR5505_CtrlChanClose(t *testing.T) {
	const N = 30
	var outer sync.WaitGroup
	outer.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer outer.Done()
			sws := newServerWatchStream()
			sws.wg.Add(2)
			go func() {
				defer func() { recover() }()
				sws.recvLoop()
			}()
			go sws.drain()
			time.Sleep(time.Microsecond)
			func() {
				defer func() { recover() }()
				sws.close()
			}()
		}()
	}
	outer.Wait()
}
