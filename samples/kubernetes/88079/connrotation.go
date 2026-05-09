package connrotation

import (
	"context"
	"net"
	"sync"
)

// Dialer is a stripped reproduction of staging/src/k8s.io/client-go/util/connrotation/connrotation.go.
// BUG (pre-PR #88079): the closeAll path reads conns map / closableConn.onClose without lock.
type DialFunc func(ctx context.Context, network, address string) (net.Conn, error)

type closableConn struct {
	net.Conn
	onClose func() // racy: written by DialContext after insertion into map
}

func (c *closableConn) Close() error {
	if c.onClose != nil {
		c.onClose()
	}
	return c.Conn.Close()
}

type Dialer struct {
	dial  DialFunc
	mu    sync.Mutex
	conns map[*closableConn]struct{}
}

func NewDialer(dial DialFunc) *Dialer {
	return &Dialer{dial: dial, conns: map[*closableConn]struct{}{}}
}

func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

func (d *Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	c, err := d.dial(ctx, network, address)
	if err != nil {
		return nil, err
	}
	cc := &closableConn{Conn: c}
	d.mu.Lock()
	d.conns[cc] = struct{}{}
	d.mu.Unlock()
	// BUG: writing onClose AFTER releasing the lock; CloseAll() reads cc.onClose.
	cc.onClose = func() {
		d.mu.Lock()
		delete(d.conns, cc)
		d.mu.Unlock()
	}
	return cc, nil
}

func (d *Dialer) CloseAll() {
	d.mu.Lock()
	conns := d.conns
	d.conns = map[*closableConn]struct{}{}
	d.mu.Unlock()
	for cc := range conns {
		// BUG: reads cc.onClose without holding any lock; race with DialContext writing it.
		_ = cc.onClose
		_ = cc.Close()
	}
}
