// Production stub for prometheus tsdb/chunkenc/bstream.go (PR #11296).
// Pre-PR: bstreamReader.loadNextBuffer reads stream slice while bstream.writeBit
// concurrently modifies the last byte of the same slice.
package chunkenc

type bstream struct {
	stream []byte
	count  uint8
}

func (b *bstream) writeBit(bit bool) {
	if b.count == 0 {
		b.stream = append(b.stream, 0)
		b.count = 8
	}
	last := len(b.stream) - 1
	b.count--
	if bit {
		b.stream[last] |= 1 << b.count // RACE: writes last byte
	}
}

type bstreamReader struct {
	stream []byte // shared slice with bstream
	pos    int
}

// loadNextBuffer reads from stream slice (RACE: stream's last byte vs writeBit).
func (r *bstreamReader) loadNextBuffer() byte {
	if r.pos >= len(r.stream) {
		return 0
	}
	v := r.stream[r.pos] // RACE
	r.pos++
	return v
}

type XORChunk struct {
	b *bstream
}

func NewXORChunk() *XORChunk {
	return &XORChunk{b: &bstream{}}
}

type xorAppender struct {
	b *bstream
}

func (a *xorAppender) Append(t int64, v float64) {
	for i := 0; i < 64; i++ {
		a.b.writeBit((uint64(v)>>uint(i))&1 == 1)
	}
	for i := 0; i < 64; i++ {
		a.b.writeBit((uint64(t)>>uint(i))&1 == 1)
	}
}

func (c *XORChunk) Appender() (*xorAppender, error) {
	return &xorAppender{b: c.b}, nil
}

type xorIterator struct {
	r *bstreamReader
}

func (i *xorIterator) Next() bool {
	if i.r.pos >= len(i.r.stream) {
		return false
	}
	_ = i.r.loadNextBuffer()
	return true
}

func (i *xorIterator) At() (int64, float64) { return 0, 0 }

func (c *XORChunk) Iterator(reuse interface{}) *xorIterator {
	return &xorIterator{r: &bstreamReader{stream: c.b.stream}}
}
