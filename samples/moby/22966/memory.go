// Production stub modeling moby PR #22966 — pkg/discovery/memory/memory.go
// data race on entries slice mutated/read without lock.
package main

type Entry struct {
	Host string
}

type Discovery struct {
	entries []*Entry
}

func NewDiscovery() *Discovery {
	return &Discovery{entries: make([]*Entry, 0)}
}

// Register appends to entries — racy with Watch (line 81 / 44 / 46 / 58 region).
func (d *Discovery) Register(host string) { // RACE write site
	d.entries = append(d.entries, &Entry{Host: host})
}

// Watch reads entries without holding a lock — racy with Register.
func (d *Discovery) Watch() int { // RACE read site
	n := 0
	for _, e := range d.entries {
		_ = e.Host
		n++
	}
	return n
}
