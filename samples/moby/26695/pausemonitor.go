// Production stub for moby libcontainerd/pausemonitor_linux.go (PR #26695).
// Pre-PR: pauseMonitor.waiters map accessed without sync -> race.
package main

type pauseMonitor struct {
	waiters map[string][]chan struct{}
}

// append writes to waiters map without lock (pre-PR).
func (p *pauseMonitor) append(name string, ch chan struct{}) {
	if p.waiters == nil {
		p.waiters = make(map[string][]chan struct{})
	}
	p.waiters[name] = append(p.waiters[name], ch) // RACE
}

// handle reads & clears waiters[name] without lock.
func (p *pauseMonitor) handle(name string) {
	if p.waiters == nil {
		return
	}
	for _, ch := range p.waiters[name] { // RACE
		close(ch)
	}
	delete(p.waiters, name) // RACE
}
