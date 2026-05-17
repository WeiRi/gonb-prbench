package grpc

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"google.golang.org/grpc/metadata"
)

// Race: BUG Invoke does `opts = append(cc.dopts.callOptions, opts...)`.
// When cc.dopts.callOptions has extra capacity, concurrent Invoke calls
// share the underlying array and race on the appended slot.
// FIX uses combine() which always allocates a new slice.
func TestRace_grpc_go_1948_invoke_callopts_append(t *testing.T) {
	noop := func(ctx context.Context, method string, req, reply interface{},
		cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error {
		return nil
	}
	cc := &ClientConn{
		dopts: dialOptions{
			unaryInt: noop,
			// extra capacity → shared underlying array on concurrent appends
			callOptions: make([]CallOption, 0, 16),
		},
	}

	var done int32
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var md metadata.MD
			opt := Header(&md)
			for j := 0; j < 5000 && atomic.LoadInt32(&done) == 0; j++ {
				_ = cc.Invoke(context.Background(), "/test/Method", nil, nil, opt)
			}
			atomic.StoreInt32(&done, 1)
		}()
	}
	wg.Wait()
}
