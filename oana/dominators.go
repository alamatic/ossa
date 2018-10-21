package oana

import (
	"github.com/alamatic/ossa"
)

// DominatorsTable is a map from each basic block to the set of basic blocks
// that are its dominators. A DominatorsTable can be constructed by calling
// FindDominators.
type DominatorsTable map[*ossa.BasicBlock]ossa.BasicBlockSet

// FindDominators calculates the dominators for the given block and all
// blocks reachable from it.
//
// Calculating dominators requires a table of predecessors provided by the
// caller. This must be the result of calling FindPredecessors with the same
// start block and no subsequent modifications to the graph beneath it, or
// the results of this function are undefined.
//
// The result is a map from each block to its dominators. Each reachable
// block must have at least one dominator: itself.
func FindDominators(start *ossa.BasicBlock, preds PredecessorsTable) DominatorsTable {
	a := dominatorsAnalyzer{
		t:     make(DominatorsTable),
		preds: preds,
	}

	ForwardDataFlow(start, a)

	return a.t
}

type dominatorsAnalyzer struct {
	t     DominatorsTable
	preds PredecessorsTable
}

func (a dominatorsAnalyzer) AnalyzeBlock(block *ossa.BasicBlock) bool {
	s, exists := a.t[block]
	if !exists {
		s = make(ossa.BasicBlockSet)
		a.t[block] = s
	}

	// Our dominator sets can only shrink as we learn more information
	// on subsequent calls, so we'll detect whether a particular block's
	// set has changed by comparing the size of the set before and after.
	priorLen := len(s)

	first := true
	for p := range a.preds[block] {
		pd, completed := a.t[p]
		if !completed {
			// Skip any predecessors that haven't had a chance to run yet.
			// This is important so we don't include these empty sets in our
			// intersection, which would cause us to preemptively remove
			// everything from our set.
			continue
		}
		if first {
			pd.AddBlocksTo(s)
			first = false
			continue
		}
		for b := range s {
			if !pd.Has(b) {
				s.Remove(b)
			}
		}
	}

	// Every block is always dominated by itself.
	s.Add(block)

	return len(s) != priorLen
}
