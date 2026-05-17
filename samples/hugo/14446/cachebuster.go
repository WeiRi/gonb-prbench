// Production stub for hugo PR #14446.
// Pre-PR: CompileConfig builds an inner closure that captures `match` from
// the enclosing scope and writes to it (instead of declaring a local with
// `:=`). Concurrent inner-closure invocations race on the captured `match`.
package main

import (
	"regexp"

	"github.com/gohugoio/hugo/common/loggers"
)

type CacheBuster struct {
	Source string
	Target string

	sourceRe *regexp.Regexp
	targetRe *regexp.Regexp
}

func (c *CacheBuster) CompileConfig(_ loggers.Logger) error {
	sre, err := regexp.Compile(c.Source)
	if err != nil {
		return err
	}
	tre, err := regexp.Compile(c.Target)
	if err != nil {
		return err
	}
	c.sourceRe = sre
	c.targetRe = tre
	return nil
}

// compiledSource returns the inner closure for matching ss against targetRe.
// BUG: captures `match` from enclosing scope and writes to it (no :=).
func (c *CacheBuster) compiledSource(src string) func(string) bool {
	if !c.sourceRe.MatchString(src) {
		return nil
	}
	var match bool
	return func(ss string) bool {
		match = c.targetRe.MatchString(ss) // RACE: write to captured `match`
		return match
	}
}
