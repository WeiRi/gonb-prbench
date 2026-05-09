package containermap

import (
	"fmt"
	"sync"
	"testing"
)

func TestRace_128657(t *testing.T) {
	cm := NewContainerMap()

	// Pre-populate
	for i := 0; i < 100; i++ {
		cid := fmt.Sprintf("container-%d", i)
		cm.Add(fmt.Sprintf("pod-%d", i%10), fmt.Sprintf("name-%d", i), cid)
	}

	var wg sync.WaitGroup
	n := 50

	// Writers: Add and Remove
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				cid := fmt.Sprintf("container-race-%d-%d", id, j)
				cm.Add(fmt.Sprintf("pod-%d", id), fmt.Sprintf("name-%d", j), cid)
				if j%2 == 0 {
					cm.RemoveByContainerID(cid)
				}
			}
		}(i)
	}

	// Readers: GetContainerRef and GetContainerID
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				cid := fmt.Sprintf("container-%d", j%100)
				cm.GetContainerRef(cid)
				cm.GetContainerID(fmt.Sprintf("pod-%d", j%10), fmt.Sprintf("name-%d", j))
			}
		}(i)
	}

	wg.Wait()
}
