package ipvs

import (
	"sync"
	"testing"
	netutils "k8s.io/utils/net"
	utilipvs "k8s.io/kubernetes/pkg/util/ipvs"
	utilipvstest "k8s.io/kubernetes/pkg/util/ipvs/testing"
)

func TestRace_107748_InPlace(t *testing.T) {
	ipvs := utilipvstest.NewFake()
	g := NewGracefulTerminationManager(ipvs)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			port := uint16(80 + id)
			for j := 0; j < 100; j++ {
				g.rsList.add(&listItem{
					VirtualServer: &utilipvs.VirtualServer{Address: netutils.ParseIPSloppy("1.1.1.1"), Protocol: "tcp", Port: port},
					RealServer:    &utilipvs.RealServer{Address: netutils.ParseIPSloppy("2.2.2.2"), Port: port},
				})
			}
		}(i)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				g.rsList.flushList(g.deleteRsFunc)
			}
		}()
	}
	wg.Wait()
}
