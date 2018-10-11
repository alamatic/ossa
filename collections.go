package ossa

// BasicBlockAdder is an interface implemented by collections that basic blocks
// can be added to, such as BasicBlockSet.
//
// This indirection is used to allow blocks to be added to a data structure
// by functions such as BasicBlock.AddSuccessors without requiring an
// incidental allocation of an intermediate set or array.
type BasicBlockAdder interface {
	Add(block *BasicBlock)
}

// basicBlockSliceBuilder is a helper to allow using the BasicBlockAdder
// interface to append to a slice in-place (by pointer).
type basicBlockSliceBuilder struct {
	s *[]*BasicBlock
}

func (b *basicBlockSliceBuilder) Add(block *BasicBlock) {
	*b.s = append(*b.s, block)
}
