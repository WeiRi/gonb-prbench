// Minimal extraction of pre-fix go-ethereum whisper/whisperv6/peer.go
// PR #16146: bloomMatch reads bloomFilter/fullNode while setBloomFilter writes them.
// Production-code paths under the upstream repo path: whisper/whisperv6/peer.go
package whisperv6

// Peer mirrors the pre-fix struct shape (subset of fields touched by race).
type Peer struct {
	trusted        bool
	powRequirement float64
	// NOTE: pre-fix had no bloomMu protecting the two fields below.
	bloomFilter []byte
	fullNode    bool
}

// Envelope stub - real type carries TopicType bloom; we only need the bloom bytes.
type Envelope struct {
	bloom []byte
}

func (e *Envelope) Bloom() []byte { return e.bloom }

// bloomFilterMatch is a faithful pre-fix port of the upstream helper.
func bloomFilterMatch(filter, sample []byte) bool {
	if filter == nil {
		return true
	}
	for i := 0; i < len(filter) && i < len(sample); i++ {
		f := filter[i]
		s := sample[i]
		if (f & s) != s {
			return false
		}
	}
	return true
}

func isFullNode(bloom []byte) bool {
	if bloom == nil {
		return true
	}
	for _, b := range bloom {
		if b != 0xff {
			return false
		}
	}
	return true
}

// bloomMatch — pre-fix version: NO mutex. PR adds peer.bloomMu.Lock/Unlock here.
// Upstream path: whisper/whisperv6/peer.go (pre-fix lines ~227-229).
func (peer *Peer) bloomMatch(env *Envelope) bool {
	return peer.fullNode || bloomFilterMatch(peer.bloomFilter, env.Bloom())
}

// setBloomFilter — pre-fix version: NO mutex. PR adds peer.bloomMu.Lock/Unlock here.
// Upstream path: whisper/whisperv6/peer.go (pre-fix lines ~232-238).
func (peer *Peer) setBloomFilter(bloom []byte) {
	peer.bloomFilter = bloom
	peer.fullNode = isFullNode(bloom)
	if peer.fullNode && peer.bloomFilter == nil {
		peer.bloomFilter = make([]byte, 64)
		for i := range peer.bloomFilter {
			peer.bloomFilter[i] = 0xff
		}
	}
}
