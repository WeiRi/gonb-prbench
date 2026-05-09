package leasing

import (
	"context"
	"sync"
)

type leaseCache struct {
	mu sync.RWMutex
}

type leasingKV struct {
	ctx      context.Context
	sessionc chan struct{}
	leases   leaseCache
}

// waitSession — BUG (pre-PR9061): reads lkv.sessionc without lock (line 449).
func (lkv *leasingKV) waitSession(ctx context.Context) error {
	c := lkv.sessionc // BUG line 449: racy read
	select {
	case <-c:
	case <-ctx.Done():
	}
	return nil
}
