package hugo7393repro

import (
	"sync"
	"testing"
)

func TestRaceContentField(t *testing.T) {
	h := &HugoSites{content: newPageMaps()}
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); h.readAndProcessContent() }()
		go func() { defer wg.Done(); h.GetContentPage() }()
	}
	wg.Wait()
}
