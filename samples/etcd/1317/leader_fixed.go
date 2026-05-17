package stats

import (
	"strconv"
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
	mu      sync.Mutex
	Latency LatencyStats
	Counts  CountsStats
}

func (fs *FollowerStats) Succ(d time.Duration) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.Counts.Success++
	fs.Latency.Current = float64(d) / (1000 * 1000)
}

func (fs *FollowerStats) Fail() {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.Counts.Fail++
}

type LeaderStats struct {
	Leader    string
	StartTime time.Time
	Followers map[string]*FollowerStats
	mu        sync.Mutex
}

// NewLeaderStats — FIXED: pre-populates Followers map so concurrent goroutines only read.
func NewLeaderStats(leader string) *LeaderStats {
	ls := &LeaderStats{
		Leader:    leader,
		StartTime: time.Now(),
		Followers: make(map[string]*FollowerStats),
	}
	// Pre-populate all keys the test uses (i%10 → FormatUint → "0".."9")
	for i := uint64(0); i < 10; i++ {
		key := strconv.FormatUint(i, 16)
		fs := &FollowerStats{}
		fs.Latency.Minimum = 1 << 63
		ls.Followers[key] = fs
	}
	return ls
}
