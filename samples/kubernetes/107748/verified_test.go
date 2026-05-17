// Race test for kubernetes-107748 (graceTerminateRSList race in flushList)
// BUG: flushList iterates q.list without lock; add() takes lock → map race.
// FIX: flushList takes lock + uses direct delete (not recursive q.remove).
package ipvs

import (
	"fmt"
	"sync"
	"testing"

	netutils "k8s.io/utils/net"
	utilipvs "k8s.io/kubernetes/pkg/util/ipvs"
)

func makeRaceListItem(i, j int) *listItem {
	vs := fmt.Sprintf("1.1.%d.%d", i, i)
	rs := fmt.Sprintf("1.1.%d.%d", i, j)
	return &listItem{
		VirtualServer: &utilipvs.VirtualServer{
			Address:  netutils.ParseIPSloppy(vs),
			Protocol: "tcp",
			Port:     uint16(80),
		},
		RealServer: &utilipvs.RealServer{
			Address: netutils.ParseIPSloppy(rs),
			Port:    uint16(80),
		},
	}
}

func TestRace_107748_FlushList(t *testing.T) {
	q := &graceTerminateRSList{
		list: make(map[string]*listItem),
	}

	// Pre-populate
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 20; j++ {
			q.add(makeRaceListItem(i, j))
		}
	}

	const N = 30
	var wg sync.WaitGroup

	// Adders: take lock
	for g := 0; g < N; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				q.add(makeRaceListItem(id+10, i))
			}
		}(g)
	}

	// FlushList callers: in BUG iterate q.list without lock → race with adders.
	noopHandler := func(rs *listItem) (bool, error) { return false, nil }
	for g := 0; g < N; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				_ = q.flushList(noopHandler)
			}
		}()
	}

	wg.Wait()
}
