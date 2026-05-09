package leastrequest

import "sync/atomic"

type scWithRPCCount struct {
	id      int
	numRPCs *int32
}

func NewSCWithCount(id int, c *int32) scWithRPCCount {
	return scWithRPCCount{id: id, numRPCs: c}
}

type DoneInfo struct{ Err error }

type DoneFunc func(DoneInfo)

type PickResult struct {
	SC   *scWithRPCCount
	Done DoneFunc
}

type Picker struct {
	scs    []scWithRPCCount
	choice int
}

func NewPicker(choice int, scs []scWithRPCCount) *Picker {
	return &Picker{choice: choice, scs: scs}
}

// Pick — BUG (pre-PR6587): plain read of *sc.numRPCs (line 45) races with atomic.AddInt32.
func (p *Picker) Pick() (PickResult, error) {
	if len(p.scs) == 0 {
		return PickResult{}, nil
	}
	picked := &p.scs[0]
	for i := range p.scs {
		sc := &p.scs[i]
		if *sc.numRPCs < *picked.numRPCs { // line 45 BUG: plain read
			picked = sc
		}
	}
	atomic.AddInt32(picked.numRPCs, 1) // atomic write
	return PickResult{SC: picked, Done: func(DoneInfo) {
		atomic.AddInt32(picked.numRPCs, -1)
	}}, nil
}
