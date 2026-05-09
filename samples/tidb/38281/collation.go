// Production stub for tidb expression/collation.go (PR #38281).
// Pre-PR: collationInfo.coerInit is plain bool. Concurrent SetCoercibility / HasCoercibility race.
package expression

type Coercibility int
type Repertoire int

type collationInfo struct {
	coer       Coercibility
	coerInit   bool // BUG: plain bool, not atomic
	repertoire Repertoire
	charset    string
	collation  string
}

func (c *collationInfo) HasCoercibility() bool {
	return c.coerInit // RACE: bare read
}

func (c *collationInfo) Coercibility() Coercibility {
	return c.coer
}

func (c *collationInfo) SetCoercibility(val Coercibility) {
	c.coer = val
	c.coerInit = true // RACE: bare write
}
