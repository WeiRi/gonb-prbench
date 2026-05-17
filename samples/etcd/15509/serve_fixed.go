package embed

import (
	"sync"
	"sync/atomic"
)

type grpcServer struct {
	id int
}

// runServe — FIXED: pass gs as parameter to goroutines instead of capturing by reference.
func runServe(wg *sync.WaitGroup) {
	defer wg.Done()
	var gs *grpcServer

	gs = &grpcServer{id: 1}
	gs1 := gs
	go func(gs *grpcServer) {
		_ = gs.id
	}(gs1)

	gs = &grpcServer{id: 2}
	gs2 := gs
	go func(gs *grpcServer) {
		_ = gs.id
	}(gs2)
	_ = atomic.AddInt32(new(int32), 1)
}
