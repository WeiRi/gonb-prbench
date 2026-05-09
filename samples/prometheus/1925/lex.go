// Production stub for prometheus PR #1925.
// Pre-PR: lex() returns an already-running lexer; caller writes
// l.seriesDesc afterwards while l.run() goroutine reads it.
package promql

type item struct {
	typ int
	val string
}

type lexer struct {
	input      string
	items      chan item
	seriesDesc bool // RACE: written by caller after lex(), read by run()
	pos        int
}

// lex starts the lexer goroutine BEFORE caller can set seriesDesc — RACE.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	// Read seriesDesc multiple times to widen the race window.
	for i := 0; i < len(l.input); i++ {
		_ = l.seriesDesc // RACE
		l.items <- item{typ: int(l.input[i]), val: string(l.input[i])}
		_ = l.seriesDesc
	}
	close(l.items)
}
