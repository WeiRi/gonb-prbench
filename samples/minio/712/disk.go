// Production stub for minio pkg/donut/disk/disk.go (PR #712).
// Pre-PR Disk has no sync.Mutex; GetFSInfo writes to fsInfo map without lock.
package disk

type Disk struct {
	path   string
	fsInfo map[string]string
}

func (d Disk) GetFSInfo() map[string]string {
	d.fsInfo["mountpoint"] = d.path
	d.fsInfo["fstype"] = "ext4"
	d.fsInfo["fsused"] = "0"
	return d.fsInfo
}
