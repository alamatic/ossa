package ossa

// BasicBlock represents a basic block in a control flow graph. A basic block
// is a straight sequence of instructions (represented as values) that always
// run as a single unit, followed by a terminator that implies zero or more
// graph edges to other successor blocks.
type BasicBlock struct {
	Instructions []*Value
	Terminator   *Terminator
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
