package apiinternal

// Reproduction of PR cockroachdb/cockroach#159666 BUG state.
// The pre-fix code declares:
//   var decoder = newDecoder()
// then later in NewAPIInternalServer mutates the decoder:
//   decoder.SetAliasTag("json")
//   decoder.IgnoreUnknownKeys(true)
// Other goroutines calling decoder.Decode() concurrently with NewAPIInternalServer
// race on the decoder's internal fields.

type Decoder struct {
	aliasTag string
	ignoreUnknown bool
}

func newDecoder() *Decoder { return &Decoder{} }

func (d *Decoder) SetAliasTag(t string)       { d.aliasTag = t }
func (d *Decoder) IgnoreUnknownKeys(b bool)   { d.ignoreUnknown = b }
func (d *Decoder) Decode() (string, bool)     { return d.aliasTag, d.ignoreUnknown }

var decoder = newDecoder()

// NewAPIInternalServer mutates the package-level decoder (BUG: race vs Decode).
func NewAPIInternalServer() {
	decoder.SetAliasTag("json")    // BUG line 27
	decoder.IgnoreUnknownKeys(true) // BUG line 28
}

// HandleRequest reads from decoder concurrently with NewAPIInternalServer (BUG).
func HandleRequest() (string, bool) {
	return decoder.Decode() // BUG line 33
}

