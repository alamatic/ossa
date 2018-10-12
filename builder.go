package ossa

// Builder is a utility for more conveniently constructing basic blocks during
// intermediate code generation in a frontend.
//
// A given builder appends instructions to a particular basic block. It is
// similar to calling the value and terminator construction functions in this
// package, but also as a side-effect appends new instructions to the basic
// block, effectively recording the order of operations.
//
// Once a terminator instruction has been appended, the builder is closed and
// any further appending calls will panic.
type Builder struct {
	block *BasicBlock
}

// NewBuilder constructs and returns a new builder.
func NewBuilder(block *BasicBlock) Builder {
	return Builder{
		block: block,
	}
}

// Block returns the block currently associated with the receiver.
func (b Builder) Block() *BasicBlock {
	return b.block
}

// SetBlock points the receiver at a different basic block. All future append
// operations will therefore apply to the new block.
func (b Builder) SetBlock(block *BasicBlock) {
	b.block = block
}

// NewBlock is a helper for allocating a new, empty basic block and wrapping
// a builder around it.
func (b Builder) NewBlock() Builder {
	block := &BasicBlock{}
	return NewBuilder(block)
}

// Open returns true if the builder is open to new instructions. That is, if
// the wrapped block does not yet have a terminator.
func (b Builder) Open() bool {
	return b.block.Terminator == nil
}

func (b Builder) appendInstruction(v *Value) *Value {
	if !b.Open() {
		panic("append to closed block")
	}
	b.block.Instructions = append(b.block.Instructions, v)
	return v
}

func (b Builder) appendTerminator(t *Terminator) *Terminator {
	if !b.Open() {
		panic("append to closed block")
	}
	b.block.Terminator = t
	return t
}

// AuxLiteral is a convenience alias for the top-level function of the
// same name. Because literals do not have side-effects, it does not append
// to the block's instruction list.
func (b Builder) AuxLiteral(v interface{}) *Value {
	return AuxLiteral(v)
}

// GlobalSym is a convenience alias for the top-level function of the
// same name. Because symbols do not have side-effects, it does not append
// to the block's instruction list.
func (b Builder) GlobalSym() *Value {
	return GlobalSym()
}

// LocalSym is a convenience alias for the top-level function of the
// same name. Because symbols do not have side-effects, it does not append
// to the block's instruction list.
func (b Builder) LocalSym() *Value {
	return LocalSym()
}

// Argument is a convenience alias for the top-level function of the
// same name. Because symbols do not have side-effects, it does not append
// to the block's instruction list.
func (b Builder) Argument() *Value {
	return Argument()
}

// Phi constructs and appends a Phi operation to the underlying block.
func (b Builder) Phi(candidates ...BasicBlockValue) *Value {
	return b.appendInstruction(Phi(candidates...))
}

// Load constructs and appends a Load operation to the underlying block.
func (b Builder) Load(ref *Value) *Value {
	return b.appendInstruction(Load(ref))
}

// Store constructs and appends a Store operation to the underlying block.
func (b Builder) Store(val, ref *Value) *Value {
	return b.appendInstruction(Store(val, ref))
}

// Call constructs and appends a Call to the underlying block.
func (b Builder) Call(callee *Value, args ...*Value) *Value {
	return b.appendInstruction(Call(callee, args...))
}

// Jump constructs a Jump terminator and uses it to terminate the underlying
// block, closing the builder.
func (b Builder) Jump(target *BasicBlock) *Terminator {
	return b.appendTerminator(Jump(target))
}

// Branch constructs a Branch terminator and uses it to terminate the underlying
// block, closing the builder.
func (b Builder) Branch(cond *Value, trueTarget, falseTarget *BasicBlock) *Terminator {
	return b.appendTerminator(Branch(cond, trueTarget, falseTarget))
}

// Switch constructs a Switch terminator and uses it to terminate the underlying
// block, closing the builder.
func (b Builder) Switch(inp *Value, defTarget *BasicBlock, cases ...BasicBlockValue) *Terminator {
	return b.appendTerminator(Switch(inp, defTarget, cases...))
}

// Return constructs a Return terminator and uses it to terminate the underlying
// block, closing the builder.
func (b Builder) Return(ret *Value) *Terminator {
	return b.appendTerminator(Return(ret))
}

// Yield constructs a Yield terminator and uses it to terminate the underlying
// block, closing the builder.
func (b Builder) Yield(resume *BasicBlock) *Terminator {
	return b.appendTerminator(Yield(resume))
}

// Await constructs a Await terminator and uses it to terminate the underlying
// block, closing the builder.
func (b Builder) Await(event *Value, resume *BasicBlock) *Terminator {
	return b.appendTerminator(Await(event, resume))
}
