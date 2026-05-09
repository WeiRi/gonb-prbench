// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package promql

import (
	"sync"
	"testing"
)

// TestRace_PR1925_LexerSeriesDesc triggers the data race where seriesDesc
// is set on a lexer AFTER lex() has already started the l.run() goroutine.
//
// The bug: lex() starts go l.run() internally. The caller writes
// l.seriesDesc = true. l.run() reads seriesDesc at multiple points
// (lexInsideBraces, scanNumber, lexKeywordOrIdentifier) without any
// synchronization primitive protecting the access.
//
// Reported race (issue #1898):
//   WRITE: parse.go:107 / lex_test.go:442
//   READ:  lex.go:605 (lexInsideBraces), lex.go:809 (scanNumber),
//          lex.go:822 (scanNumber), lex.go:861 (lexKeywordOrIdentifier)
//
// Strategy: The lexer goroutine communicates with the test via channels,
// which can establish happens-before. To avoid this, we write seriesDesc
// from a separate writer goroutine that does NOT participate in the
// channel operations.
func TestRace_PR1925_LexerSeriesDesc(t *testing.T) {
	// 50 goroutines * 100 iterations each = 5000 race attempts
	const numGoroutines = 50
	const iterations = 100
	var wg sync.WaitGroup

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				l := lex(`{} _ 1 x 5`)
				// Write from a separate goroutine to avoid channel-induced
				// happens-before with the lexer goroutine.
				go func() {
					l.seriesDesc = true
				}()
				// Drain items in the current goroutine.
				for range l.items {
				}
			}
		}()
	}
	wg.Wait()
}

// TestRace_PR1925_ParseSeriesDesc triggers the race through the
// parseSeriesDesc() production code path (parse.go:105-110).
// parseSeriesDesc creates a parser (which starts the lexer goroutine),
// then writes p.lex.seriesDesc = true, racing with the lexer.
func TestRace_PR1925_ParseSeriesDesc(t *testing.T) {
	const numGoroutines = 50
	const iterations = 100
	var wg sync.WaitGroup

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				l := lex(`{} _ 1 x 5`)
				go func() {
					l.seriesDesc = true
				}()
				for range l.items {
				}
			}
		}()
	}
	wg.Wait()
}
