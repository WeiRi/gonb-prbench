// Same-package whitebox race test for grpc-go#5927
// PR diff file: xds/internal/xdsclient/clientimpl_loadreport.go
//
// Pre-fix bug: ReportLoad() unlocks authorityMu BEFORE calling a.refLocked().
// refCount is then mutated by 2 goroutines: one from ReportLoad's a.refLocked()
// (no lock), another from c.unrefAuthority() (with lock).
//
// To get a frame in clientimpl_loadreport.go, we call c.ReportLoad() directly.
// We pre-seed c.authorities[cfg.String()] so newAuthorityLocked short-circuits
// and never touches transport — avoiding gRPC dial setup.
package xdsclient

import (
	"sync"
	"testing"

	"ase/grpc-go-5927/cache"
	"ase/grpc-go-5927/bootstrap"
)

func TestRace_5927(t *testing.T) {
	cfg := &bootstrap.ServerConfig{ServerURI: "race-target"}
	cfgStr := cfg.String()

	c := &clientImpl{
		authorities:     map[string]*authority{cfgStr: {refCount: 0, serverCfg: cfg}},
		idleAuthorities: cache.NewTimeoutCache(0),
	}

	const N = 30
	const ITERS = 50
	var wg sync.WaitGroup
	wg.Add(N + N)

	// Goroutine A: call ReportLoad which on pre-fix unlocks BEFORE refLocked().
	// a.reportLoad() then nil-derefs transport — recover to swallow.
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			defer func() { _ = recover() }()
			for j := 0; j < ITERS; j++ {
				_, _ = c.ReportLoad(cfg)
			}
		}()
	}

	// Goroutine B: call c.unrefAuthority(a) which decrements refCount under lock.
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			defer func() { _ = recover() }()
			a := c.authorities[cfgStr]
			for j := 0; j < ITERS; j++ {
				c.authorityMu.Lock()
				if a.refCount > 0 {
					a.refCount--
				} else {
					a.refCount++
				}
				c.authorityMu.Unlock()
			}
		}()
	}
	wg.Wait()
}
