package clientv3

import (
	"context"
	"sync"
)

// FIXED: MDLocal wrapped in struct with shared sync.Mutex pointer.
// Pointer ensures copies share the same mutex.
type MDLocal struct {
	mu *sync.Mutex
	m  map[string][]string
}

func PairsLocal(kvs ...string) MDLocal {
	md := MDLocal{mu: new(sync.Mutex), m: map[string][]string{}}
	for i := 0; i < len(kvs); i += 2 {
		md.m[kvs[i]] = append(md.m[kvs[i]], kvs[i+1])
	}
	return md
}

func (m MDLocal) Set(k, v string) {
	m.mu.Lock()
	m.m[k] = []string{v}
	m.mu.Unlock()
}

func (m MDLocal) Get(k string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.m[k]
}

type ctxKeyMD struct{}

func NewOutgoingContextLocal(ctx context.Context, md MDLocal) context.Context {
	return context.WithValue(ctx, ctxKeyMD{}, md)
}

func FromOutgoingContextLocal(ctx context.Context) (MDLocal, bool) {
	v, ok := ctx.Value(ctxKeyMD{}).(MDLocal)
	return v, ok
}
