package tcpproxy

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestRace_etcd_7361_runmonitor_loop_var(t *testing.T) {
	tp := &TCPProxy{
		MonitorInterval: time.Microsecond,
		donec:           make(chan struct{}),
	}
	for _, addr := range []string{"127.0.0.1:1", "127.0.0.1:2", "127.0.0.1:3"} {
		tp.remotes = append(tp.remotes, &remote{addr: addr, inactive: true})
	}

	go tp.runMonitor()

	var done int32
	time.AfterFunc(2*time.Second, func() { atomic.StoreInt32(&done, 1) })
	for atomic.LoadInt32(&done) == 0 {
		time.Sleep(10 * time.Millisecond)
	}
}
