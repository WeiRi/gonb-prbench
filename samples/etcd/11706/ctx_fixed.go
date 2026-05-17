package clientv3

import "context"

// WithRequireLeader — FIXED: uses locked access via MDLocal.Get
func WithRequireLeader(ctx context.Context) context.Context {
	md, ok := FromOutgoingContextLocal(ctx)
	if !ok {
		return ctx
	}
	// FIXED: use locked Get instead of direct iteration
	_ = md.Get("key1")
	return NewOutgoingContextLocal(ctx, md)
}
