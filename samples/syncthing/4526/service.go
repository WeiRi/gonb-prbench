// Race-trigger test for syncthing-4526; see README.md for usage.

package connections

import "sync"

type internalConn struct{ priority int }
type DeviceID [32]byte

type dialTarget struct {
	priority int
	uri      string
	deviceID DeviceID
}

func (t dialTarget) Dial() (internalConn, error) {
	// reads multiple fields of t (line 70 BUG)
	return internalConn{priority: t.priority}, nil
}

// DialParallel mirrors the racy inner loop where goroutine captures tgt.
func DialParallel(deviceID DeviceID, tgts []dialTarget) (internalConn, bool) {
	var wg sync.WaitGroup
	res := make(chan internalConn, len(tgts))
	for _, tgt := range tgts { // line 67 BUG: range overwrites tgt
		wg.Add(1)
		go func() { // line 69 BUG: captures tgt
			conn, err := tgt.Dial() // line 70 BUG: reads tgt fields
			if err == nil {
				res <- conn
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(res)
	c, ok := <-res
	return c, ok
}
