package clientv3

import "sync"

type Logger struct{ name string }

func NewNopLogger() *Logger { return &Logger{} }

type Client struct {
	lg   *Logger
	lgMu *sync.RWMutex
}

// SetLogger — BUG (pre-PR20192): writes c.lg without holding c.lgMu.
func (c *Client) SetLogger(lg *Logger) {
	c.lg = lg // line 16 BUG
}

// GetLogger — BUG: reads c.lg without holding c.lgMu.
func (c *Client) GetLogger() *Logger {
	return c.lg // line 21 BUG
}
