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
	ID         string
	State      StateType
	LeaderInfo leaderInfo
	sync.Mutex
}

// SetState — BUG (pre-PR5157): writes ss.State + LeaderInfo without lock.
func (ss *ServerStats) SetState(s StateType, name string) {
	ss.State = s                    // line 27 BUG
	ss.LeaderInfo.Name = name       // line 28 BUG
	ss.LeaderInfo.StartTime = time.Now() // line 29 BUG
}

// GetState — BUG: reads without lock.
func (ss *ServerStats) GetState() (StateType, string, time.Time) {
	return ss.State, ss.LeaderInfo.Name, ss.LeaderInfo.StartTime // line 34 BUG
}
