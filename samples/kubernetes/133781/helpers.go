package v1

// Stripped reproduction of pkg/apis/scheduling/v1/helpers.go pre-PR #133781.
// BUG: SystemPriorityClasses returns a shared (package-global) slice; concurrent readers
// and writers (callers mutating the returned PriorityClasses) race on the elements' fields.
// The PR returns a fresh slice / deep-copy.

type PriorityClass struct {
	Kind  string
	Value int32
	Name  string
}

var systemPriorityClasses = []*PriorityClass{
	{Name: "system-node-critical", Value: 2000001000, Kind: ""},
	{Name: "system-cluster-critical", Value: 2000000000, Kind: ""},
}

// SystemPriorityClasses — BUG: returns the shared slice.
func SystemPriorityClasses() []*PriorityClass {
	return systemPriorityClasses
}

// IsKnownSystemPriorityClass reads Name + Value + Kind on the shared elements.
func IsKnownSystemPriorityClass(name string, value int32, _ bool) bool {
	for _, c := range systemPriorityClasses {
		// reads Kind => races with TestRace_133781 writing c.Kind = "PriorityClass"
		_ = c.Kind
		if c.Name == name && c.Value == value {
			return true
		}
	}
	return false
}
