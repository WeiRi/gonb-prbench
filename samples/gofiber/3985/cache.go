// Production stub for gofiber/middleware/cache (PR #3985).
// Models the racy `item` struct: updateData writes body+expire concurrently
// while readData/readDate read them without sync.
package cache

type item struct {
	body   []byte
	expire int64
	status int
	headers map[string][]string
	cType  []byte
	cEnc   []byte
}

// updateData mirrors the writer path that PR #3985 fixed (was unsynchronized).
func (i *item) updateData(body []byte, expire int64) {
	i.body = body
	i.expire = expire
}

// readData reads the body field without synchronization (race vs updateData).
func (i *item) readData() []byte {
	return i.body
}

// readDate reads the expire field without synchronization (race vs updateData).
func (i *item) readDate() int64 {
	return i.expire
}
