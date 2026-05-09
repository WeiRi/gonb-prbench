// Production stub modeling moby PR #21677 — layer/ro_layer.go data race
// on referenceCount field (incremented from layer_store.go without lock).
package main

type roLayer struct {
	chainID         string
	referenceCount  int
}

type layerStore struct {
	layerMap map[string]*roLayer
}

func NewLayerStore() *layerStore {
	return &layerStore{layerMap: make(map[string]*roLayer)}
}

func (ls *layerStore) register(id string) *roLayer {
	l, ok := ls.layerMap[id]
	if !ok {
		l = &roLayer{chainID: id}
		ls.layerMap[id] = l
	}
	return l
}

// hold is the racy mutator (mirrors layer_store.go:337 reference bump).
func (l *roLayer) hold() { // RACE write site
	l.referenceCount++
}

// refCount is the racy reader (mirrors ro_layer.go:86 reference read).
func (l *roLayer) refCount() int { // RACE read site
	return l.referenceCount
}
