package oana

import (
	"github.com/alamatic/ossa"
)

// NaturalLoop represents a natural loop discovered within a control flow
// graph. A natural loop is defined by a "back edge", which is an edge in
// the graph from a node back to another node that dominates it.
type NaturalLoop struct {
	// Head and Tail describe the back edge that define the loop, which is
	// from Tail to Head.
	Head, Tail *ossa.BasicBlock
}

// FindNaturalLoops uses the given dominators table to detect any natural
// loops, appending each one found to the given slice "to" which
// may be nil.
//
// The caller must provide the result of calling FindDominators with some
// start block, without any modification to the graph in the mean time, or
// the result is undefined.
func FindNaturalLoops(doms DominatorsTable, to []NaturalLoop) []NaturalLoop {
	for block, blockDoms := range doms {
		// If any of the successors of our block also dominate it then
		// we have found a loop.
		block.AddSuccessors(basicBlockAdderFunc(func(succ *ossa.BasicBlock) {
			if !blockDoms.Has(succ) {
				return
			}

			to = append(to, NaturalLoop{
				Head: succ,
				Tail: block,
			})
		}))

	}
	return to
}

// FindBody finds the set of basic blocks that form the body of the receiving
// loop, which includes the loop's head and tail as well as any ancestors of
// tail that are not also ancestors of head.
//
// The caller must provide the result of calling FindPredecessors with the
// same start block that was used to produce the dominators table that the
// loop was detected from, without any modification to the graph in the mean
// time, or the result is undefined.
func (l *NaturalLoop) FindBody(preds PredecessorsTable) ossa.BasicBlockSet {
	ret := ossa.NewBasicBlockSet(l.Head)
	q := newBlockLIFO(4)
	q.Add(l.Tail)
	for !q.Empty() {
		block := q.Next()
		if !ret.Has(block) {
			ret.Add(block)
			for pb := range preds[block] {
				q.Add(pb)
			}
		}
	}
	return ret
}
