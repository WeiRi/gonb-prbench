// Production-code stand-in for hugo helpers/general.go (pre-fix DistinctLogger
// path). Kept in a separate file so race detector reports frames here, not in
// the _test.go.
package hugo6410repro

import (
	"fmt"
	"strings"
	"sync"
)

type LogPrinter interface {
	Println(a ...interface{})
}

type stubLogger struct{}

func (stubLogger) Println(a ...interface{}) {}

type DistinctLogger struct {
	sync.RWMutex
	logger LogPrinter
	m      map[string]bool
}

func (l *DistinctLogger) Println(v ...interface{}) {
	logStatement := strings.TrimSpace(fmt.Sprintln(v...))
	l.print(logStatement)
}

func (l *DistinctLogger) print(logStatement string) {
	l.RLock()
	if l.m[logStatement] {
		l.RUnlock()
		return
	}
	l.RUnlock()
	l.Lock()
	if !l.m[logStatement] {
		l.logger.Println(logStatement)
		l.m[logStatement] = true
	}
	l.Unlock()
}

func NewDistinctErrorLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: stubLogger{}}
}

var DistinctErrorLog = NewDistinctErrorLogger()

func InitLoggers() {
	DistinctErrorLog = NewDistinctErrorLogger()
}
