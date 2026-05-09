// BUG version of k8s.io/client-go/tools/cache/fifo.go (pre-PR-90825).
// closed is guarded by a *separate* closedLock, while cond uses f.lock.
// Close() does Broadcast() OUTSIDE f.lock, so a Pop() that is between
// "cond.Wait() returned" and "IsClosed() returns true" may have ALREADY
// passed IsClosed(false) check and re-entered Wait() — and never gets the
// next broadcast. -> Pop blocks forever after Close().
package pock90825

import (
	"errors"
	"sync"
)

var ErrFIFOClosed = errors.New("FIFO closed")

type FIFO struct {
	lock sync.Mutex
	cond sync.Cond

	items map[string]struct{}
	queue []string

	closed     bool
	closedLock sync.Mutex
}

func NewFIFO() *FIFO {
	f := &FIFO{items: map[string]struct{}{}}
	f.cond.L = &f.lock
	return f
}

func (f *FIFO) Add(key string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if _, ok := f.items[key]; !ok {
		f.queue = append(f.queue, key)
		f.items[key] = struct{}{}
	}
	f.cond.Broadcast()
}

func (f *FIFO) Close() {
	f.closedLock.Lock()
	defer f.closedLock.Unlock()
	f.closed = true
	f.cond.Broadcast() // BUG: broadcast not under f.lock
}

func (f *FIFO) IsClosed() bool {
	f.closedLock.Lock()
	defer f.closedLock.Unlock()
	return f.closed
}

type PopProcessFunc func(string) error

// Pop — BUG: checks IsClosed() AFTER cond.Wait, but closedLock is independent
// so the broadcast on Close() may slip through unnoticed.
func (f *FIFO) Pop(process PopProcessFunc) (string, error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	for {
		for len(f.queue) == 0 {
			if f.IsClosed() {
				return "", ErrFIFOClosed
			}
			f.cond.Wait()
		}
		k := f.queue[0]
		f.queue = f.queue[1:]
		delete(f.items, k)
		_ = process(k)
		return k, nil
	}
}
