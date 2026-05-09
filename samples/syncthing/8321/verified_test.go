package connections

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRace_8321(t *testing.T) {
	s := &service{
		conns:  make(chan internalConn, 100),
		hellos: make(chan *connWithHello, 1000),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// drain
	go func() {
		for {
			select {
			case <-s.hellos:
			case <-ctx.Done():
				return
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = s.HandleConns(ctx)
	}()

	for i := 0; i < 5000 && ctx.Err() == nil; i++ {
		select {
		case s.conns <- internalConn{id: i, state: i}:
		case <-ctx.Done():
			break
		}
	}
	<-ctx.Done()
	wg.Wait()
}
