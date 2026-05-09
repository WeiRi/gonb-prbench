package v3rpc

import "sync/atomic"

var ProgressReportIntervalMilliseconds int32

func sendLoop() int32 {
	return ProgressReportIntervalMilliseconds // non-atomic read RACE
}

func SetInterval(v int32) {
	atomic.StoreInt32(&ProgressReportIntervalMilliseconds, v)
}
