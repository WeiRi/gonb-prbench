// Pre-fix pacer/pacer.go from rclone PR #140.
// BUG: SetMinSleep / SetMaxSleep / SetRetries write fields concurrently with
// endCall / Call which read them — no mutex.
package pacer

import "time"

type Pacer struct {
	minSleep      time.Duration
	maxSleep      time.Duration
	decayConstant uint
	pacer         chan struct{}
	sleepTime     time.Duration
	retries       int
}

func New() *Pacer {
	return &Pacer{minSleep: 10 * time.Millisecond, maxSleep: 2 * time.Second, sleepTime: 10 * time.Millisecond, retries: 3, pacer: make(chan struct{}, 1)}
}

// SetMinSleep — pacer.go:44 pre-fix (no lock).
func (p *Pacer) SetMinSleep(t time.Duration) *Pacer {
	p.minSleep = t
	p.sleepTime = p.minSleep
	return p
}

// SetMaxSleep — pacer.go:50.
func (p *Pacer) SetMaxSleep(t time.Duration) *Pacer {
	p.maxSleep = t
	p.sleepTime = p.minSleep
	return p
}

// SetDecayConstant — pacer.go:62.
func (p *Pacer) SetDecayConstant(d uint) *Pacer {
	p.decayConstant = d
	return p
}

// SetRetries — pacer.go:69.
func (p *Pacer) SetRetries(r int) *Pacer {
	p.retries = r
	return p
}

// endCall — pacer.go:95 reads/writes sleepTime without lock.
func (p *Pacer) endCall(again bool) {
	oldSleepTime := p.sleepTime
	if again {
		p.sleepTime *= 2
		if p.sleepTime > p.maxSleep {
			p.sleepTime = p.maxSleep
		}
	} else {
		_ = oldSleepTime
	}
}

// Call — pacer.go:142 reads p.retries without lock.
func (p *Pacer) Call() int {
	return p.retries
}
