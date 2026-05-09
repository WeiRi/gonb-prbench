package rafthttp

import (
	"sync"
	"testing"

	
)

func TestRace_10826(t *testing.T) {
	tr := &Transport{
		peers: map[ID]Peer{},
	}
	for i := 1; i <= 8; i++ {
		tr.peers[ID(i)] = newFakePeer()
	}

	const N = 8
	const ITERS = 500
	var wgW sync.WaitGroup
	var wgR sync.WaitGroup
	stop := make(chan struct{})

	wgW.Add(N)
	for i := 0; i < N; i++ {
		go func(seed int) {
			defer wgW.Done()
			for j := 0; j < ITERS; j++ {
				tr.mu.Lock()
				id := ID(100 + seed*100 + j%32)
				if _, ok := tr.peers[id]; ok {
					delete(tr.peers, id)
				} else {
					tr.peers[id] = newFakePeer()
				}
				tr.mu.Unlock()
			}
		}(i)
	}

	wgR.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wgR.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				tr.Pause()
				tr.Resume()
			}
		}()
	}

	wgW.Wait()
	close(stop)
	wgR.Wait()
}
