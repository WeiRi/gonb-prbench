package admission

import "sync"

// Reproduction of PR cockroachdb/cockroach#131292 / 131293
// "admission: lock work queue before reading waiting length"

type heap struct {
	items []int
}
func (h *heap) Len() int { return len(h.items) }
func (h *heap) Push(x int) { h.items = append(h.items, x) }

type tenantInfo struct {
	waitingWorkHeap *heap
}

type WorkQueue struct {
	mu     sync.Mutex
	tenant *tenantInfo
}

// Admit reads tenant.waitingWorkHeap.Len() WITHOUT holding q.mu (BUG).
func (q *WorkQueue) Admit() int {
	return q.tenant.waitingWorkHeap.Len() // BUG (line 26)
}

// granted dequeues; reads heap len under no lock (BUG).
func (q *WorkQueue) Granted() int {
	return q.tenant.waitingWorkHeap.Len() // BUG (line 31)
}

// Push appends to heap UNDER lock (concurrent writer).
func (q *WorkQueue) Push(x int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tenant.waitingWorkHeap.Push(x)
}

func NewWorkQueue() *WorkQueue {
	return &WorkQueue{tenant: &tenantInfo{waitingWorkHeap: &heap{}}}
}

