// Production stub for gofiber/fiber ctx.go (PR #4062).
// Pre-PR: DefaultCtx.Value reads c.fasthttp without sync; release() sets it nil.
package fiber

import "sync"

type fakeFastHTTPCtx struct {
	mu     sync.Mutex
	values map[any]any
}

func (f *fakeFastHTTPCtx) UserValue(key any) any {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.values[key]
}

func newFakeFastHTTPCtx() *fakeFastHTTPCtx {
	return &fakeFastHTTPCtx{values: map[any]any{"k": "v"}}
}

type DefaultCtx struct {
	fasthttp *fakeFastHTTPCtx
}

func NewDefaultCtx() *DefaultCtx {
	return &DefaultCtx{fasthttp: newFakeFastHTTPCtx()}
}

// Value reads c.fasthttp without sync (BUG).
func (c *DefaultCtx) Value(key any) any {
	return c.fasthttp.UserValue(key) // RACE
}

// release writes c.fasthttp = nil without sync (BUG).
func (c *DefaultCtx) release() {
	c.fasthttp = nil // RACE
}
