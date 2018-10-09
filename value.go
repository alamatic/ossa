package ossa

// Value is the most fundamental type in ossa, representing a node in the SSA
// graph.
type Value struct {
	op Op

	// args is zero or more argument values whose meaning depends on the
	// indicated operation.
	args []*Value

	// aux is an auxillary native Go value
	aux interface{}

	// For ops that use three or fewer args, this can be used as the backing
	// array for args, avoiding another allocation. The size 3 is chosen
	// to make just enough room for call instructions that are representing
	// either unary or binary operators (where the first element is a
	// representation of the operator itself.)
	argsBuf [3]*Value
}

var Void *Value

func (v *Value) Op() Op {
	return v.op
}

// AuxLiteral constructs a new Value with OpAuxLiteral.
func AuxLiteral(v interface{}) *Value {
	return &Value{
		op:  OpAuxLiteral,
		aux: v,
	}
}

// GlobalSym constructs a new global symbol. A global symbol's value pointer
// its identity; it contains no further data.
func GlobalSym() *Value {
	return &Value{
		op: OpGlobalSym,
	}
}

// LocalSym constructs a new local symbol. A local symbol's value pointer
// its identity; it contains no further data.
func LocalSym() *Value {
	return &Value{
		op: OpLocalSym,
	}
}

// Argument constructs a new argument placeholder. An argument's value pointer
// its identity; it contains no further data.
func Argument() *Value {
	return &Value{
		op: OpArgument,
	}
}

// Phi constructs a Phi node, representing the join of various possible source
// values at the entry into a basic block.
func Phi(candidates ...BasicBlockValue) *Value {
	return &Value{
		op:   OpPhi,
		args: bbvsAsArgs(candidates),
	}
}

// Load constructs a Load instruction value, reading from the memory object
// described by the given value.
func Load(ref *Value) *Value {
	v := &Value{
		op: OpLoad,
	}
	v.args = v.argsBuf[:1]
	v.args[0] = ref
	return v
}

// Store constructs a Store instruction value, writing the given value to the
// the memory object described by the given ref value.
func Store(val *Value, ref *Value) *Value {
	v := &Value{
		op: OpStore,
	}
	v.args = v.argsBuf[:2]
	v.args[0] = v
	v.args[1] = ref
	return v
}

// Call constructs a Call instruction value, which represents calling the
// callee value with the given argument values.
//
// This instruction type can be used to represent calls to both user-defined
// functions and fundamental operations in a language, with the former perhaps
// represented by global symbols while the latter may be represented by
// AuxLiteral values that would not be otherwise representable in the language.
func Call(callee *Value, args ...*Value) *Value {
	v := &Value{
		op: OpCall,
	}
	aa := v.bufForArgs(len(args) + 1)
	aa = append(aa, callee)
	for _, a := range args {
		aa = append(aa, a)
	}
	v.args = aa
	return v
}

// bufForArgs returns a zero-length value slice with at least the given capacity
// that can be used as the arguments for the receiving value.
//
// bufForArgs either allocates a slice with the given capacity ready to have
// arguments appended to it, or returns a slice backed by v.argsBuf if
// its length is enough to contain the args with no further allocation.
func (v *Value) bufForArgs(capacity int) []*Value {
	if len(v.argsBuf) >= capacity {
		return v.argsBuf[:0]
	}
	return make([]*Value, 0, capacity)
}
