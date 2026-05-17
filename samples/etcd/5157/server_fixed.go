package stats

import (
	"sync"
	"time"
)

type StateType int

const (
	StateFollower StateType = iota
	StateLeader
)

type leaderInfo struct {
	Name      string
	StartTime time.Time
}

type ServerStats struct {
	mu         sync.Mutex
	ID         string
	State      StateType
	LeaderInfo leaderInfo
	sync.Mutex
}

// SetState — BUG (pre-PR5157): writes ss.State + LeaderInfo without lock.
func (ss *ServerStats) SetState(s StateType, name string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.State = s                         // line 27 BUG
	ss.LeaderInfo.Name = name            // line 28 BUG
	ss.LeaderInfo.StartTime = time.Now() // line 29 BUG
}

// GetState — BUG: reads without lock.
func (ss *ServerStats) GetState() (StateType, string, time.Time) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.State, ss.LeaderInfo.Name, ss.LeaderInfo.StartTime // line 34 BUG
}
