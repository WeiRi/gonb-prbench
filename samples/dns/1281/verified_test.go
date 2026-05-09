package dns

import (
	"sync"
	"testing"
	"time"
)

// TestRace_dns1281_Client_Dialer: concurrent ExchangeContext calls write
// shared c.Dialer pointer while another goroutine reads via Dial — races
// fire under -race.
func TestRace_dns1281_Client_Dialer(t *testing.T) {
	c := &Client{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			c.ExchangeContext(time.Duration(i) * time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = c.Dial()
		}
	}()
	wg.Wait()
}
