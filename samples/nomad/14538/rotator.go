// Production stub for nomad client/logmon/logging/rotator.go (PR #14538).
// Pre-PR: oldestLogFileIdx field accessed concurrently between purgeOldFiles
// (writer) and nextFile (reader) without a mutex.
package logging

type FileRotator struct {
	path             string
	prefix           string
	maxFiles         int
	fileSize         int64
	logger           interface{}
	logFileIdx       int
	oldestLogFileIdx int
}

func NewFileRotator(path, prefix string, maxFiles int, fileSize int64, logger interface{}) (*FileRotator, error) {
	return &FileRotator{
		path: path, prefix: prefix, maxFiles: maxFiles, fileSize: fileSize, logger: logger,
	}, nil
}

func (f *FileRotator) Close() error { return nil }

// purgeOldFiles writes oldestLogFileIdx without lock (line 313 upstream).
func (f *FileRotator) purgeOldFiles(idx int) {
	f.oldestLogFileIdx = idx
}

// nextFile reads oldestLogFileIdx without lock (line 181 upstream).
func (f *FileRotator) nextFile() int {
	return f.logFileIdx - f.oldestLogFileIdx
}
