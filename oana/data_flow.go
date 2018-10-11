package oana

import (
	"github.com/alamatic/ossa"
)

// BlockAnalyzer is the interface implemented by types to be used with
// block-oriented analysis algorithms.
//
// For example, implementations of this type are used to drive data flow
// analyses using both ForwardDataFlow and BackwardDataFlow.
type BlockAnalyzer interface {
	// AnalyzeBlock is called for each block visited by an algorithm, in
	// an order defined by that algorithm.
	//
	// The implementer should update its analysis data structures to
	// incorporate the given block and return true if and only if the effective
	// result was changed by those updates. For most algorithms the same
	// block may be passed multiple times, e.g. if there are loops in the
	// control flow graph. On subsequent calls, information about other
	// predecessor blocks may have changed. It is required that the
	// implementation eventually reach a fixpoint, such that all future
	// calls to AnalyzeBlock for any block in the graph will return false.
	//
	// There is no mechanism for an implementer to directly return an error,
	// so instead an implementation must record any errors as part of its
	// result state and then keep returning false from AnalyzeBlock until
	// no further calls arrive, before allowing the algorithm caller to
	// inspect the result through implementation-specific mechanisms.
	AnalyzeBlock(block *ossa.BasicBlock) (changed bool)
}

// BlockAnalyzerFunc is an implementation of BlockAnalyzer that calls a
// single callback function with the same signature as AnalyzeBlock.
type BlockAnalyzerFunc func(block *ossa.BasicBlock) (changed bool)

func (f BlockAnalyzerFunc) AnalyzeBlock(block *ossa.BasicBlock) bool {
	return f(block)
}

// ForwardDataFlow performs a forward data flow analysis on the control flow
// graph entered at the given start block, driven by the given analyzer
// implementation.
//
// The analyzer will first be called with the start block. If it returns true
// then each of the successors of that block will be added to a work queue
// and called in turn, with the result processed in the same way. In most
// uses of this function, the first call to the analyzer for a given block
// will return true to indicate that the data for that block was initialized,
// but it is also valid to return false in order to skip processing successors
// altogether, if for example an error occurs or enough information has already
// been gathered.
//
// Note that it is not guaranteed that all of a block's predecessors will be
// called before that block, since that is not possible in general in the
// presence of loops. Analyzers must be prepared to tolerate incomplete
// information and expect to visit the same block again later once more
// predecessors have produced data.
//
// The ordering of visiting blocks will be consistent for a particular version
// of this module, but the ordering is not part of the function's contract and
// may change in future versions.
func ForwardDataFlow(start *ossa.BasicBlock, analyzer BlockAnalyzer) {
	q := newBlockLIFO(6) // enough capacity to process a flat-ish CFG without further allocation
	q.Add(start)

	for !q.Empty() {
		block := q.Next()
		changed := analyzer.AnalyzeBlock(block)
		if changed {
			// Add all successors to the processing queue.
			l := q.Length()
			block.AddSuccessors(q)

			// We prefer to visit successors in the reverse order to what
			// AddSuccessors generates, because in the usual form of loops
			// this allows us to analyze the loop body and then re-analyze
			// the loop header before moving on to the block after the loop,
			// thus only visiting that final block once rather than twice.
			q.ReverseTopN(q.Length() - l)
		}
	}
}
