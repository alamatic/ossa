package ossa

// BasicBlock represents a basic block in a control flow graph. A basic block
// is a straight sequence of instructions (represented as values) that always
// run as a single unit, followed by a terminator that implies zero or more
// graph edges to other successor blocks.
type BasicBlock struct {
	Instructions []*Value
	Terminator   *Terminator
}

// AddSuccessors adds the successors of this block to the given set, modifying
// it in-place.
func (b *BasicBlock) AddSuccessors(to BasicBlockSet) {
	b.Terminator.AddSuccessors(to)
}

// AddReachable adds to the given set all of the blocks that are reachable
// from the receiver, including the receiver itself. The set is modified
// in-place.
//
// This method assumes that the given set, if not empty, was built by a prior
// call to AddReachable, and so any blocks already present in the set is
// already accompanied by all of the blocks reachable from it. In other words,
// it will not visit any of the descendents of any blocks already present in
// the set, even if they are reachable from the receiver.
func (b *BasicBlock) AddReachable(to BasicBlockSet) {
	// todo serves as a work queue, but we'll use stack discipline for it just
	// to keep things simple, since our order of processing is not important.
	todo := make([]*BasicBlock, 0, 4) // cap 4 just to give us some room for simple graphs without more allocation
	todo = append(todo, b)
	for len(todo) > 0 {
		block := todo[len(todo)-1]
		todo = todo[:len(todo)-1]
		if to.Has(block) {
			continue
		}
		to.Add(block)
		todo = b.Terminator.AppendSuccessors(todo)
	}
}

// BasicBlockValue represents a (BasicBlock, Value) pair, used in a small
// number of value factory functions.
type BasicBlockValue struct {
	Block *BasicBlock
	Value *Value
}

func bbvsAsArgs(pairs []BasicBlockValue) []*Value {
	args := make([]*Value, 0, len(pairs)*2)
	for _, pair := range pairs {
		args = append(args, &Value{
			op:  opBasicBlock,
			aux: pair.Block,
		})
		args = append(args, pair.Value)
	}
	return args
}
