// Production stub for hugo markup/goldmark/hugocontext/hugocontext.go (PR #12413).
// Pre-PR Wrap returns buf.Bytes() from a sync.Pool *bytes.Buffer; the returned
// slice aliases the pooled buffer and races with concurrent re-uses.
package hugocontext

import (
	"bytes"
	"strconv"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// Wrap mirrors the racy pre-PR signature returning []byte from a pooled buffer.
func Wrap(payload []byte, ord uint64) []byte {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.WriteString("<!--ctx[")
	buf.WriteString(strconv.FormatUint(ord, 10))
	buf.WriteString("]-->")
	buf.Write(payload)
	out := buf.Bytes() // RACE: aliases pool buffer, returned to callers concurrently
	bufPool.Put(buf)
	return out
}
