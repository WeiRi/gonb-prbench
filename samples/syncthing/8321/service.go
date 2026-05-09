// Race-trigger test for syncthing-8321; see README.md for usage.

package connections

import "context"

type internalConn struct {
	id    int
	state int
}

type connWithHello struct {
	c      internalConn
	remote int
}

type service struct {
	conns  chan internalConn
	hellos chan *connWithHello
}

// HandleConns mirrors the racy do() function from the BUG diff.
func (s *service) HandleConns(ctx context.Context) error {
	var c internalConn // BUG: declared OUTSIDE loop
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c = <-s.conns: // line 34 BUG: overwrites loop var
		}
		go func() { // captures c
			select {
			case s.hellos <- &connWithHello{c: c, remote: c.id}: // line 41 BUG: reads c
			case <-ctx.Done():
			}
		}()
	}
}
