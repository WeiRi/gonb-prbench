// Production stub modeling etcd PR #10772 — mvcc/index.go racy treeIndex.
// Mirrors the bug: Put/Get on internal map without proper locking.
package main

type keyIndex struct {
	rev int64
}

type treeIndex struct {
	tree map[string]*keyIndex
}

func newTreeIndex() *treeIndex {
	return &treeIndex{tree: make(map[string]*keyIndex)}
}

// Put writes into the tree map without taking a lock — racy with Get.
func (ti *treeIndex) Put(key string, rev int64) { // RACE write site
	ki, ok := ti.tree[key]
	if !ok {
		ki = &keyIndex{}
		ti.tree[key] = ki
	}
	ki.rev = rev
}

// Get reads from the tree map without taking a lock — racy with Put.
func (ti *treeIndex) Get(key string) int64 { // RACE read site
	ki, ok := ti.tree[key]
	if !ok {
		return 0
	}
	return ki.rev
}
