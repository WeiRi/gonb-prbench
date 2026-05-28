// Regression test for kubernetes#98653
// Bug: StreamWatcher.receive() sends Error to result channel without
// coordinating with Stop(). When Stop() runs after receive's stopping() check
// but before the send, the result channel has no receiver → send blocks forever.
// Fix: a `done` channel signaled by Stop() lets the send select-default out.
//
// Oracle: PANIC (goroutine leak). Dump all goroutine stacks before t.Fatalf
// so generate_code_candidates can parse production frames (streamwatcher.go).

package watch

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type sw98653Decoder struct {
	emit chan error
}

func (d *sw98653Decoder) Decode() (EventType, apiruntime.Object, error) {
	err, ok := <-d.emit
	if !ok {
		return "", nil, errors.New("decoder closed")
	}
	return "", nil, err
}
func (d *sw98653Decoder) Close() { close(d.emit) }

type sw98653Reporter struct{}

func (sw98653Reporter) AsObject(err error) apiruntime.Object {
	return &sw98653Obj{msg: err.Error()}
}

type sw98653Obj struct{ msg string }

func (o *sw98653Obj) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }
func (o *sw98653Obj) DeepCopyObject() apiruntime.Object { return &sw98653Obj{msg: o.msg} }

func TestStreamWatcherRace_98653(t *testing.T) {
	const iters = 30
	const perIter = 3

	startGoroutines := runtime.NumGoroutine()

	for it := 0; it < iters; it++ {
		for k := 0; k < perIter; k++ {
			d := &sw98653Decoder{emit: make(chan error, 1)}
			r := sw98653Reporter{}
			sw := NewStreamWatcher(d, r)
			d.emit <- errors.New("synthetic-decode-error-98653")
			time.Sleep(2 * time.Millisecond)
			sw.Stop()
		}
		time.Sleep(10 * time.Millisecond)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		cur := runtime.NumGoroutine()
		if cur-startGoroutines <= 2 {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	leak := runtime.NumGoroutine() - startGoroutines

	// Dump all live goroutine stacks so race_report.txt contains parseable
	// frames pointing at streamwatcher.go (the leaked receive goroutines).
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, true)
	fmt.Fprintln(os.Stderr, "goroutine dump (PANIC oracle 98653, leak):")
	fmt.Fprintln(os.Stderr, string(buf[:n]))
	t.Fatalf("PANIC oracle fired: %d leaked goroutines after Stop (bug 98653)", leak)
}
