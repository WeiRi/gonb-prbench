package status

import (
	"fmt"
	"sync"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/klog/v2/ktesting"
)

// TestRace_136205_StartTime_Sharing triggers the data race where
// updateStatusInternal shares the StartTime pointer from the cached
// status with the new status, and normalizeStatus later writes to
// that shared pointer, racing with concurrent readers via GetPodStatus.
func TestRace_136205_StartTime_Sharing(t *testing.T) {
	logger, _ := ktesting.NewTestContext(t)
	syncer := newTestManager(&fake.Clientset{})

	// Set initial pod status with a non-nil StartTime so the code path
	// in updateStatusInternal that copies the pointer is exercised.
	pod := getTestPod()
	startTime := metav1.Now()
	initialStatus := v1.PodStatus{
		StartTime: &startTime,
		Conditions: []v1.PodCondition{
			{
				Type:   v1.PodReady,
				Status: v1.ConditionTrue,
			},
		},
	}
	syncer.SetPodStatus(logger, pod, initialStatus)

	numWriters := 25
	numReaders := 25
	iterations := 200
	var wg sync.WaitGroup

	// Writer goroutines: repeatedly call SetPodStatus, which calls
	// updateStatusInternal -> shares StartTime pointer -> normalizeStatus
	// writes to *StartTime via normalizeTimeStamp.
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				now := metav1.Now()
				status := v1.PodStatus{
					StartTime: &now,
					Conditions: []v1.PodCondition{
						{
							Type:   v1.PodReady,
							Status: v1.ConditionTrue,
						},
					},
				}
				syncer.SetPodStatus(logger, pod, status)
			}
		}(i)
	}

	// Reader goroutines: repeatedly call GetPodStatus and access
	// StartTime fields, racing with normalizeStatus writing to the
	// shared StartTime pointer (before the DeepCopy fix).
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				status, ok := syncer.GetPodStatus(pod.UID)
				if ok && status.StartTime != nil {
					// Access StartTime fields to race with normalizeStatus writes
					_ = status.StartTime.Time
					_ = status.StartTime.Format(time.RFC3339)
					_ = fmt.Sprintf("%v", status.StartTime)
				}
			}
		}(i)
	}

	wg.Wait()
}
