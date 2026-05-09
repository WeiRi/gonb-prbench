package vsphere

import "testing"

func TestRace_90348_LoopVarCapture(t *testing.T) {
	const N = 500
	dcNodes := map[string][]string{
		"dc1": {"n1", "n2", "n3"},
		"dc2": {"n4", "n5", "n6"},
		"dc3": {"n7", "n8", "n9"},
		"dc4": {"n10", "n11", "n12"},
	}
	for i := 0; i < N; i++ {
		ProcessDisksAreAttached(dcNodes)
	}
}
