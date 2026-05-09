package stats

import (
	"sync"
	"time"
)

type LatencyStats struct {
	Current           float64
	Average           float64
	averageSquare     float64
	StandardDeviation float64
	Minimum           float64
	Maximum           float64
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

// BUG (pre-PR1316): no lock around concurrent stat mutation.
func (fs *FollowerStats) Succ(d time.Duration) {
	total := float64(fs.Counts.Success)*fs.Latency.Average + float64(d)/(1000*1000)
	fs.Counts.Success++
	fs.Latency.Average = total / float64(fs.Counts.Success)
	fs.Latency.Current = float64(d) / (1000 * 1000)
	if fs.Latency.Current > fs.Latency.Maximum {
		fs.Latency.Maximum = fs.Latency.Current
	}
}

func (fs *FollowerStats) Fail() {
	fs.Counts.Fail++
}
