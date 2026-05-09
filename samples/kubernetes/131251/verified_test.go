// Regression test for kubernetes#131251
// Bug: realImageGCManager.freeImage() called `delete(im.imageRecords, key)`
// WITHOUT acquiring im.imageRecordsLock. Concurrent readers (e.g., detectImages,
// imageRecordsLen) hold the lock — race detector fires.
// Fix: Lock/Unlock around the delete.
package images

import (
	"context"
	"sync"
	"testing"

	"k8s.io/kubernetes/pkg/kubelet/container"
	statstest "k8s.io/kubernetes/pkg/kubelet/server/stats/testing"
)

// TestImageGCManager_FreeImage_LockOnDelete_131251 reproduces the data race
// fixed by PR 131251. Bug version: races; Fix version: clean.
func TestImageGCManager_FreeImage_LockOnDelete_131251(t *testing.T) {
	mockStatsProvider := new(statstest.MockProvider)
	im, fakeRuntime := newRealImageGCManager(ImageGCPolicy{}, mockStatsProvider)

	// Pre-populate imageRecords + fakeRuntime image list
	const N = 100
	for i := 0; i < N; i++ {
		key := fmtKey(i)
		im.imageRecords[key] = &imageRecord{}
		fakeRuntime.ImageList = append(fakeRuntime.ImageList, container.Image{ID: key})
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Reader: takes lock and reads imageRecords (mirrors detectImages-like access)
	go func() {
		defer wg.Done()
		for i := 0; i < N*5; i++ {
			_ = im.imageRecordsLen()
		}
	}()

	// Writer: calls freeImage which deletes from imageRecords.
	// In BUG state this delete runs WITHOUT the lock → race with reader.
	// In FIX state freeImage acquires the lock → no race.
	go func() {
		defer wg.Done()
		ctx := context.Background()
		for i := 0; i < N; i++ {
			key := fmtKey(i)
			ev := evictionInfo{id: key, imageRecord: imageRecord{}}
			_ = im.freeImage(ctx, ev, "test")
		}
	}()
	wg.Wait()
}

func fmtKey(i int) string {
	return "img-" + intToStr(i)
}

func intToStr(i int) string {
	if i == 0 {
		return "0"
	}
	digits := []byte{}
	for i > 0 {
		digits = append([]byte{byte('0' + i%10)}, digits...)
		i /= 10
	}
	return string(digits)
}
