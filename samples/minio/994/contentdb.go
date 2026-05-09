// Production stub for minio contentdb/contentdb.go (PR #994).
// Pre-PR Init() has no mutex; multiple goroutines racing on isInitialized + extDB map.
package contentdb

import "errors"

var (
	isInitialized bool
	extDB         map[string]string
)

func Init() error {
	if isInitialized {
		// Even when "already inited", the buggy code still writes extDB without sync.
		extDB = make(map[string]string)
		extDB["txt"] = "text/plain"
		isInitialized = true
		return nil
	}
	if err := loadDB(); err != nil {
		return err
	}
	extDB["txt"] = "text/plain"
	extDB["json"] = "application/json"
	isInitialized = true
	return nil
}

func loadDB() error {
	if extDB == nil {
		extDB = make(map[string]string)
	}
	return nil
}

func Lookup(ext string) (string, bool) {
	v, ok := extDB[ext]
	return v, ok
}

func MustLookup(ext string) string {
	v, ok := Lookup(ext)
	if !ok {
		panic(errors.New("not found"))
	}
	return v
}
