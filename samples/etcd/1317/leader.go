package stats

import (
	"sync"
	"time"
)

type LatencyStats struct {
	Average           float64
	averageSquare     float64
	StandardDeviation float64
	Minimum           float64
	Maximum           float64
	Current           float64
}

type CountsStats struct {
	Fail    uint64
	Success uint64
}

type FollowerStats struct {
	Latency LatencyStats
	Counts  CountsStats
	sync.Mutex
}

func (fs *FollowerStats) Succ(d time.Duration) {
	fs.Counts.Success++
	fs.Latency.Current = float64(d) / (1000 * 1000)
}

func (fs *FollowerStats) Fail() {
	fs.Counts.Fail++
}

type LeaderStats struct {
	Leader    string
	StartTime time.Time
	// BUG (pre-PR1317): direct map access without lock. Fix introduced sync.Mutex.
	Followers map[string]*FollowerStats
	sync.Mutex
}

func NewLeaderStats(leader string) *LeaderStats {
	return &LeaderStats{
		Leader:    leader,
		StartTime: time.Now(),
		Followers: make(map[string]*FollowerStats),
	}
}
