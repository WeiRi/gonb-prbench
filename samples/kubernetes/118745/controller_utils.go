package controller

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Stripped reproduction of pkg/controller/controller_utils.go pre-PR #118745.
// BUG: logExpectations uses %#v which reflects ControlleeExpectations's add/del fields
// without atomic.Load, racing with concurrent atomic.AddInt64.

type ControlleeExpectations struct {
	add int64
	del int64
	key string
}

func (e *ControlleeExpectations) Add(addAdds, addDels int64) {
	atomic.AddInt64(&e.add, addAdds)
	atomic.AddInt64(&e.del, addDels)
}

// logExpectations — BUG: %#v reflectively reads add/del without atomic.Load.
func logExpectations(w io.Writer, prefix string, exp *ControlleeExpectations) {
	fmt.Fprintf(w, "%s %#v\n", prefix, exp)
}
