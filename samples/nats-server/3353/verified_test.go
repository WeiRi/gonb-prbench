// VERIFIED race reproducer for nats-server PR #3353
// "[FIXED] Data race"
// https://github.com/nats-io/nats-server/pull/3353
//
// Original racy file: server/filestore.go (populateGlobalPerSubjectInfo
// ~line 5018 pre-fix).
// Pre-fix calls mb.readPerSubjectInfo(false) without acquiring mb.mu;
// concurrent writers mutate mb fields under mb.mu.Lock.
// Fix: take mb.mu.Lock() / defer mb.mu.Unlock() around the call.
//
// 4 writers x 8 readers x 5000 ops, golang:1.21 -race.
package buggy

import (
	"fmt"
	"sync"
	"testing"
)

func TestRacePopulateGlobalPerSubjectInfo(t *testing.T) {
	const N = 5000
	const W = 4
	const R = 8
	mb := &msgBlock{}
	fs := &fileStore{}
	var wg sync.WaitGroup
	for i := 0; i < W; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < N; j++ {
				mb.writePerSubject(fmt.Sprintf("k%d-%d", id, j), j)
			}
		}(i)
	}
	for i := 0; i < R; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				fs.populateGlobalPerSubjectInfo(mb)
			}
		}()
	}
	wg.Wait()
}
