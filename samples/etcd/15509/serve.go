package embed

import (
	"sync"
	"sync/atomic"
)

type grpcServer struct {
	id int
}

// runServe — BUG (pre-PR15509): outer-scope `gs` reassigned across two paths,
// goroutines reference outer var → reads race writes from sibling reassign.
func runServe(wg *sync.WaitGroup) {
	defer wg.Done()
	var gs *grpcServer

	gs = &grpcServer{id: 1} // line 37 first assign
	go func() {
		_ = gs.id // racy read of outer-scope gs
	}()

	gs = &grpcServer{id: 2} // line 40 second assign races with the goroutine read
	go func() {
		_ = gs.id
	}()
	_ = atomic.AddInt32(new(int32), 1)
}
