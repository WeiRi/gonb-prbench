// Production stub for minio memory/intelligent.go (PR #568).
// Pre-PR ExpireObjects iterates r.items without holding the lock.
package memory

import (
	"sync"
	"time"
)

type Intelligent struct {
	lock     sync.Mutex
	items    map[string][]byte
	created  map[string]time.Time
	maxSize  int64
	expireIn time.Duration
}

func NewIntelligent(maxSize int64, expireIn time.Duration) *Intelligent {
	return &Intelligent{
		items:    make(map[string][]byte),
		created:  make(map[string]time.Time),
		maxSize:  maxSize,
		expireIn: expireIn,
	}
}

func (r *Intelligent) Set(key string, value []byte) {
	r.lock.Lock()
	r.items[key] = value
	r.created[key] = time.Now()
	r.lock.Unlock()
}

// ExpireObjects iterates r.items WITHOUT holding the lock (pre-PR bug).
func (r *Intelligent) ExpireObjects(maxAge time.Duration) {
	now := time.Now()
	for k, t := range r.created { // RACE: concurrent map iter w/ Set's writes
		if now.Sub(t) > maxAge {
			delete(r.items, k)
			delete(r.created, k)
		}
	}
}
