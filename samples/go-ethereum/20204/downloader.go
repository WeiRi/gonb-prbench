// Minimal extraction of pre-fix go-ethereum eth/downloader/downloader.go
// for PR #20204. Pre-fix race: an anonymous goroutine closure captures the
// outer `stateSync` variable; the outer scope later reassigns
// `stateSync = d.syncState(...)`, so the goroutine's read of stateSync races
// with the outer write. PR fix replaces the closure with a function-typed
// `closeOnErr := func(s *stateSync) { ... }` invoked as `go closeOnErr(sync)`.
// Production-code path: eth/downloader/downloader.go
package downloader

import "sync"

type stateSync struct {
	root int
	mu   sync.Mutex
	err  error
}

func (s *stateSync) Wait() error { return nil }

type Downloader struct{}

func (d *Downloader) syncState(root int) *stateSync {
	return &stateSync{root: root}
}

// processFastSyncContent — pre-fix simplified loop reproducing the variable-capture race.
// Upstream path: eth/downloader/downloader.go (pre-fix line ~1577 anon goroutine).
func (d *Downloader) processFastSyncContent(iters int) {
	stateSync := d.syncState(0)
	// Pre-fix: anonymous goroutine captures `stateSync` by reference.
	go func() {
		_ = stateSync.Wait()
		// Read the captured stateSync field — races with reassignment below.
		_ = stateSync.root
	}()
	for i := 1; i < iters; i++ {
		// This reassignment is the racy write captured by the closure above.
		stateSync = d.syncState(i)
		go func() {
			_ = stateSync.Wait()
			_ = stateSync.root
		}()
	}
	_ = stateSync
}
