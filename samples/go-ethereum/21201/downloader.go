// PR #21201 - eth/downloader/downloader.go - data race on Downloader.mode
// (SyncMode). Pre-fix: d.mode is a plain SyncMode field written by
// synchronise() and read by Progress, syncWithPeer, fetchHeight, findAncestors,
// fetchHeaders, processHeaders. PR fix: convert to uint32 + atomic ops via
// getMode() / atomic.StoreUint32.
// Production-code path: eth/downloader/downloader.go (pre-fix line ~103, 419, 261, 471, 642, 724, 749, 807, 880, 1017).
package downloader

type SyncMode uint32

const (
	FullSync  SyncMode = 0
	FastSync  SyncMode = 1
	LightSync SyncMode = 2
)

type Downloader struct {
	mode SyncMode // PRE-FIX: plain field, no sync.
}

func NewDownloader() *Downloader {
	return &Downloader{mode: FullSync}
}

// synchronise — pre-fix writer of d.mode.
// Upstream: eth/downloader/downloader.go (pre-fix line ~419).
func (d *Downloader) Synchronise(mode SyncMode) {
	d.mode = mode
}

// Progress — pre-fix reader of d.mode (no atomic).
// Upstream: eth/downloader/downloader.go (pre-fix line ~261).
func (d *Downloader) Progress() SyncMode {
	return d.mode
}
