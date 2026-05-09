package client

// Client is a stand-in for moby client/client.go (BUG state).
// Bug: customHTTPHeaders map is read on every API call while being
// concurrently written by SetCustomHTTPHeaders, without synchronization.

type Client struct {
	customHTTPHeaders map[string]string // BUG: unsynchronized
}

func NewClient() *Client {
	return &Client{customHTTPHeaders: map[string]string{}}
}

// SetCustomHTTPHeaders mutates the map (race write).
func (c *Client) SetCustomHTTPHeaders(h map[string]string) {
	c.customHTTPHeaders = h
}

// AddHeader writes to the map (race write).
func (c *Client) AddHeader(k, v string) {
	c.customHTTPHeaders[k] = v
}

// GetHeader reads the map (race read).
func (c *Client) GetHeader(k string) string {
	return c.customHTTPHeaders[k]
}

// CountHeaders iterates the map (race iterator vs writer).
func (c *Client) CountHeaders() int {
	n := 0
	for range c.customHTTPHeaders {
		n++
	}
	return n
}
