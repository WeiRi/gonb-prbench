// In-place race test for etcd-7070: package=lease, uses upstream Lease.
// Bug: lessor.go — Demote() writes Lease.expiry, Remaining() reads it,
// both without atomic. Race on bare int64 (monotime.Time).
// PR fix: use atomic.StoreUint64/LoadUint64 for expiry access.
package lease

import (
	"sync"
	"testing"

	monotime "github.com/coreos/etcd/pkg/monotime"
)

func TestRace_7070_InPlace(t *testing.T) {
	const N = 50
	const ITERS = 200

	for trial := 0; trial < ITERS; trial++ {
		l := &Lease{ID: LeaseID(1)}
		var wg sync.WaitGroup
		wg.Add(N * 2)

		for i := 0; i < N; i++ {
			// Half goroutines: call Demote -> write l.expiry (lessor.go:342)
			go func() {
				defer wg.Done()
				l.expiry = monotime.Now()
			}()
			// Half goroutines: call Remaining -> read l.expiry (lessor.go Remaining())
			go func() {
				defer wg.Done()
				_ = l.Remaining()
			}()
		}
		wg.Wait()
	}
}
