package clientv3

import "context"

// WithRequireLeader — BUG (pre-PR11706): reads md directly without copying;
// concurrent md.Set() callers race with the iteration here.
func WithRequireLeader(ctx context.Context) context.Context {
	md, ok := FromOutgoingContextLocal(ctx)
	if !ok {
		return ctx
	}
	// BUG: read md fields directly (would race with concurrent md.Set())
	for k, vs := range md {
		_ = k
		_ = vs
	}
	return NewOutgoingContextLocal(ctx, md)
}
