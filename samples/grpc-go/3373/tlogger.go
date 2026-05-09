package grpctest

import "regexp"

type TLogger struct {
	errors map[*regexp.Regexp]int // BUG: map without lock
}

func NewTLogger() *TLogger {
	return &TLogger{errors: map[*regexp.Regexp]int{}}
}

func (g *TLogger) AddExpect(re *regexp.Regexp, n int) {
	g.errors[re] = n // line 20 write
}

func (g *TLogger) expected(msg string) bool {
	for re := range g.errors { // line 22 iterate - races with AddExpect
		if re.MatchString(msg) {
			return true
		}
	}
	return false
}

func (g *TLogger) Update() {
	for re := range g.errors {
		_ = re
	}
}

func (g *TLogger) EndTest() {
	for re := range g.errors {
		delete(g.errors, re)
	}
}
