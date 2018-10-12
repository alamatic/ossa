package oana

import (
	"github.com/alamatic/ossa"
)

// AllPredecessors calculates the predecessors for the given block and all
// blocks reachable from it, by inverting all of the "successor" edges
// implied by the block terminators.
//
// The result is a map from each block to its predecessors. Each reachable
// block must have at least one predecessor by definition, since otherwise
// it would not be reachable.
func AllPredecessors(start *ossa.BasicBlock) map[*ossa.BasicBlock]ossa.BasicBlockSet {
	ret := make(map[*ossa.BasicBlock]ossa.BasicBlockSet)
	seen := make(ossa.BasicBlockSet)

	q := newBlockLIFO(6)
	q.Add(start)
	for !q.Empty() {
		pred := q.Next()
		seen.Add(pred)
		pred.AddSuccessors(basicBlockAdderFunc(func(succ *ossa.BasicBlock) {
			if _, exists := ret[succ]; !exists {
				ret[succ] = make(ossa.BasicBlockSet)
			}
			ret[succ].Add(pred)
			if !seen.Has(succ) {
				q.Add(succ)
			}
		}))
	}

	return ret
}

// basicBlockAdderFunc is a bit of a cheat to let us use functions that take
// basicBlockAdderFuncs as a mapping function over whatever blocks are
// added.
type basicBlockAdderFunc func(block *ossa.BasicBlock)

func (f basicBlockAdderFunc) Add(block *ossa.BasicBlock) {
	f(block)
}
