package grpc

import "sync"

type ccResolverWrapper struct {
	mu       sync.Mutex
	resolver *struct{ id int }
}

// NewCCResolverWrapper — BUG (pre-PR3090): assigns ccr.resolver while a goroutine
// in another path reads ccr.resolver concurrently.
func NewCCResolverWrapper(stop chan struct{}) *ccResolverWrapper {
	ccr := &ccResolverWrapper{}
	go func() {
		<-stop
		_ = ccr.resolver // line 25 BUG racy read
	}()
	ccr.resolver = &struct{ id int }{id: 1} // line 56 write without lock
	return ccr
}
