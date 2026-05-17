package server

import "sync"

type sourceInfo struct {
	wg  sync.WaitGroup
	qch chan struct{}
}

// RacyWaitGroupWait calls wg.Wait() — races with concurrent Add/Done
func (si *sourceInfo) RacyWaitGroupWait() {
	si.wg.Wait()
}

// RacyWaitGroupAddDone calls Add(1) then Done() — races with concurrent Wait
func (si *sourceInfo) RacyWaitGroupAddDone() {
	si.wg.Add(1)
	si.wg.Done()
}
