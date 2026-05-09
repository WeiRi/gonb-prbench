package grpc

import (
	"errors"
	"time"
)

var (
	minConnectTimeout = 20 * time.Second
	errCancel         = errors.New("cancel")
)

func resetTransport(stop <-chan struct{}) error {
	for {
		select {
		case <-stop:
			return errCancel
		default:
		}
		dialDuration := minConnectTimeout // line 28 BUG plain read
		_ = dialDuration
	}
}
