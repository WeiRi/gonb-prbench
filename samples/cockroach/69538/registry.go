package metric

// Registry is a simplified reproducer of cockroach pkg/util/metric/registry.go
// (BUG state): metrics map accessed concurrently from AddMetric and Each
// without locking.
type Registry struct {
	metrics map[string]int64 // BUG: unsynchronized concurrent access
}

func NewRegistry() *Registry {
	return &Registry{metrics: map[string]int64{}}
}

// AddMetric: registry.go:72 area — concurrent map write.
func (r *Registry) AddMetric(name string, val int64) {
	r.metrics[name] = val
}

// Each: registry.go:143 area — concurrent map read iteration.
func (r *Registry) Each(fn func(name string, val int64)) {
	for k, v := range r.metrics { // race read
		fn(k, v)
	}
}

// Get: registry.go:161 area — concurrent map read.
func (r *Registry) Get(name string) int64 {
	return r.metrics[name]
}
