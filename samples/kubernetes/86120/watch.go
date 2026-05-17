package watch
import "sync"
type Shared struct { mu sync.Mutex; val int64 }
func New() *Shared { return &Shared{} }
func (s *Shared) Write(v int64) { s.val = v }
func (s *Shared) Read() int64 { return s.val }
