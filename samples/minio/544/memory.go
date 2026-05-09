// Production stub for minio pkg/storage/drivers/memory/memory.go (PR #544).
// Pre-PR expireObjects reads len(memory.objectMetadata) without lock.
package memory

import (
	"io"
	"sync"
	"time"
)

type memoryDriver struct {
	lock           sync.Mutex
	buckets        map[string]string
	objectMetadata map[string][]byte
	maxSize        int64
	expireDur      time.Duration
}

type Driver interface {
	CreateBucket(name, acl string) error
	CreateObject(bucket, key, contentType, etag string, data io.Reader) error
}

func Start(maxSize int64, expireDur time.Duration) (chan<- struct{}, <-chan error, Driver) {
	d := &memoryDriver{
		buckets:        make(map[string]string),
		objectMetadata: make(map[string][]byte),
		maxSize:        maxSize,
		expireDur:      expireDur,
	}
	go d.expireObjects()
	return nil, nil, d
}

func (d *memoryDriver) CreateBucket(name, acl string) error {
	d.lock.Lock()
	d.buckets[name] = acl
	d.lock.Unlock()
	return nil
}

func (d *memoryDriver) CreateObject(bucket, key, ctype, etag string, data io.Reader) error {
	buf := make([]byte, 0, 64)
	d.lock.Lock()
	d.objectMetadata[bucket+"/"+key] = buf
	d.lock.Unlock()
	return nil
}

// expireObjects reads len(d.objectMetadata) WITHOUT holding the lock (pre-PR).
func (d *memoryDriver) expireObjects() {
	for {
		// RACE: read len of map vs concurrent CreateObject writes (under lock)
		_ = len(d.objectMetadata)
		time.Sleep(time.Microsecond)
	}
}
