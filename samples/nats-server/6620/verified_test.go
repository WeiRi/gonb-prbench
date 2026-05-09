// VERIFIED race reproducer for nats-server PR #6620
// "[FIXED] Data race in configureAccounts"
// https://github.com/nats-io/nats-server/pull/6620
//
// Original racy file: server/server.go (configureAccounts ~line 1383-1388)
// Pre-fix: s.sys.account is read AFTER s.mu.Unlock(); fix snapshots sysAcc before.
//
// Recipe: 8 readers x 8 writers x 5000 ops, golang:1.21 -race => race.
package buggy

import (
	"fmt"
	"sync"
	"testing"
)

func TestRaceConfigureAccountsSysAccount(t *testing.T) {
	s := &Server{sys: &internal{account: &account{name: "init"}}}
	const W = 8
	const R = 8
	const N = 5000
	var wg sync.WaitGroup
	for i := 0; i < W; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < N; j++ {
				s.rotateSysAccount(fmt.Sprintf("a%d-%d", id, j))
			}
		}(i)
	}
	for i := 0; i < R; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				_ = s.configureAccounts()
			}
		}()
	}
	wg.Wait()
}
