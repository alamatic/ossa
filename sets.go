package ossa

// BasicBlockSet is a data structure for a set of basic blocks.
type BasicBlockSet map[*BasicBlock]struct{}

// NewBasicBlockSet is a helper for constructing a basic block set with an
// initial set of members. It is also valid to construct an empty set with
// the "make" function.
func NewBasicBlockSet(blocks ...*BasicBlock) BasicBlockSet {
	ret := make(BasicBlockSet)
	for _, b := range blocks {
		ret.Add(b)
	}
	return ret
}

// Has returns true only if the given block is in the set.
func (s BasicBlockSet) Has(block *BasicBlock) bool {
	_, ok := s[block]
	return ok
}

// Add inserts the given block into the set. It is a no-op if the block is
// already present in the set.
func (s BasicBlockSet) Add(block *BasicBlock) {
	s[block] = struct{}{}
}

// Remove removes the given block from the set. It is a no-op if the block is
// not already in the set.
func (s BasicBlockSet) Remove(block *BasicBlock) {
	delete(s, block)
}

// RemoveAll removes all members from the set, making the set empty.
func (s BasicBlockSet) RemoveAll() {
	for block := range s {
		delete(s, block)
	}
}

// AppendBlocks appends to the given slice all of the blocks in the set (in
// a non-deterministic order) and returns the new slice.
func (s BasicBlockSet) AppendBlocks(to []*BasicBlock) []*BasicBlock {
	if len(s) == 0 {
		return to
	}
	needCap := len(to) + len(s)
	if cap(to) < needCap {
		newCapTo := cap(to) * 2
		newCapFrom := len(s) * 2 // always at least 2, because we return early if len=0
		var new []*BasicBlock
		if newCapFrom < newCapTo {
			new = make([]*BasicBlock, len(to), newCapFrom)
		} else {
			new = make([]*BasicBlock, len(to), newCapTo)
		}
		copy(new, to)
		to = new
	}
	for v := range s {
		to = append(to, v)
	}
	return to
}

// AddBlocksTo adds to the given adder all of the blocks in the set in a
// non-deterministic order.
func (s BasicBlockSet) AddBlocksTo(to BasicBlockAdder) {
	for v := range s {
		to.Add(v)
	}
}

// ValueSet is a data structure for a set of values.
type ValueSet map[*Value]struct{}

// Has returns true only if the given value is in the set.
func (s ValueSet) Has(value *Value) bool {
	_, ok := s[value]
	return ok
}

// Add inserts the given value into the set. It is a no-op if the value is
// already present in the set.
func (s ValueSet) Add(value *Value) {
	s[value] = struct{}{}
}

// Remove removes the given value from the set. It is a no-op if the value is
// not already in the set.
func (s ValueSet) Remove(value *Value) {
	delete(s, value)
}

// OpSet is a data structure for a set of opcodes.
type OpSet map[*Op]struct{}

// Has returns true only if the given opcode is in the set.
func (s OpSet) Has(op *Op) bool {
	_, ok := s[op]
	return ok
}

// Add inserts the given opcode into the set. It is a no-op if the opcode is
// already present in the set.
func (s OpSet) Add(op *Op) {
	s[op] = struct{}{}
}

// Remove removes the given opcode from the set. It is a no-op if the opcode is
// not already in the set.
func (s OpSet) Remove(op *Op) {
	delete(s, op)
}
