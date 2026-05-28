// Companion to watch_based_manager.go — defines cacheStore (referenced by
// the verified_test.go). BUG: initialized is a plain bool; reads racy with
// writes done under cacheStore.lock.
package manager

import (
	"sync"

	"k8s.io/client-go/tools/cache"
)

type cacheStore struct {
	cache.Store
	initialized bool
	lock        sync.Mutex
}
