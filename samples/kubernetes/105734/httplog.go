// Pre-fix httplog.go from PR #105734 (apiserver/server/httplog).
// BUG: respLogger.Addf, AddKeyValue, Log access addedInfo / addedKeyValuePairs
// without a mutex. Concurrent goroutines (timeout filter + auth filter) call
// these → data race on the slice / strings.Builder fields.
package httplog

import (
	"strings"
)

type respLogger struct {
	addedInfo          strings.Builder
	addedKeyValuePairs []interface{}
}

func newRespLogger() *respLogger {
	return &respLogger{addedKeyValuePairs: make([]interface{}, 0, 16)}
}

// Addf — httplog.go:195 in pre-fix (no lock).
func (rl *respLogger) Addf(format string) {
	rl.addedInfo.WriteString("\n")
	rl.addedInfo.WriteString(format)
}

// AddKeyValue — httplog.go:206 in pre-fix.
func (rl *respLogger) AddKeyValue(key string, value interface{}) {
	rl.addedKeyValuePairs = append(rl.addedKeyValuePairs, key, value)
}

// Log — httplog.go reads addedKeyValuePairs without lock.
func (rl *respLogger) Log() int {
	keysAndValues := []interface{}{"verb", "GET"}
	keysAndValues = append(keysAndValues, rl.addedKeyValuePairs...) // racy READ of slice header
	_ = rl.addedInfo.String()
	return len(keysAndValues)
}
