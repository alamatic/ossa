package ossa

import (
	"fmt"
)

// Terminator represents the edge between one basic block and zero or more
// other successor basic blocks in a control flow graph.
type Terminator struct {
	op Op

	// args is zero or more argument values whose meaning depends on the
	// indicated operation. Some elements of this slice may not use both
	// fields of struct BasicBlockValue, depending on the needs of the op.
	args []BasicBlockValue

	// For ops that use two or fewer args, this can be used as the backing
	// array for args, avoiding another allocation. The size 3 is chosen
	// to make just enough room for call instructions that are representing
	// either unary or binary operators (where the first element is a
	// representation of the operator itself.)
	argsBuf [2]BasicBlockValue
}

// Jump constructs an unconditional jump terminator leading to the given
// other basic block.
func Jump(target *BasicBlock) *Terminator {
	t := &Terminator{
		op: OpJump,
	}
	t.argsBuf[0].Block = target
	t.args = t.argsBuf[:1]
	return t
}

// Branch constructs a conditional branch terminator with the given condition
// value and pair of target basic blocks.
func Branch(cond *Value, trueTarget, falseTarget *BasicBlock) *Terminator {
	t := &Terminator{
		op: OpBranch,
	}
	t.argsBuf[0].Value = cond
	t.argsBuf[0].Block = trueTarget
	t.argsBuf[1].Block = falseTarget // argsBuf[1].Value is unused
	t.args = t.argsBuf[:2]
	return t
}

// Switch constructs a conditional switch terminator with the given input
// value, default target basic block, and zero or more conditional branch
// pairs.
func Switch(inp *Value, defTarget *BasicBlock, cases ...BasicBlockValue) *Terminator {
	t := &Terminator{
		op: OpSwitch,
	}
	aa := t.bufForArgs(len(cases) + 1)
	aa = append(aa, BasicBlockValue{
		Value: inp,
		Block: defTarget,
	})
	aa = append(aa, cases...)
	t.args = aa
	return t
}

// Return constructs a terminator that exits the current function with the
// given return value. This terminator produces no successors.
func Return(ret *Value) *Terminator {
	t := &Terminator{
		op: OpReturn,
	}
	t.argsBuf[0].Value = ret
	t.args = t.argsBuf[:1]
	return t
}

// Yield constructs a terminator that acts as a yield point for coroutines.
// Yield indicates that the routine wishes to yield control to another routine.
// The exact behavior of a yield is ultimately decided by the language runtime;
// for languages that don't use coroutines, do not generate Yield terminators.
//
// The given basic block is the point where execution will continue after the
// coroutine is resumed.
func Yield(resume *BasicBlock) *Terminator {
	t := &Terminator{
		op: OpYield,
	}
	t.argsBuf[0].Block = resume
	t.args = t.argsBuf[:1]
	return t
}

// Await constructs a terminator that acts as an async blocking point for
// coroutines. This is similar to Await except also takes an argument for
// some language-defined event value (promise, etc) that must occur or complete
// before the routine can resume.
//
// The given basic block is the point where execution will continue after the
// coroutine is resumed.
func Await(event *Value, resume *BasicBlock) *Terminator {
	t := &Terminator{
		op: OpAwait,
	}
	t.argsBuf[0].Value = event
	t.argsBuf[0].Block = resume
	t.args = t.argsBuf[:1]
	return t
}

// Unreachable is a special terminator that has no behavior and no successors.
// This should be used only in situations where the language frontend can
// guarantee control can never reach a certain point (or it would be undefined
// behavior to do so).
//
// For example, this might be emitted immediately after a call instruction which
// the frontend knows cannot actually return in practice, e.g. if it exits the
// program, or just blocks/loops forever.
//
// Although this is a variable, callers are forbidden from assigning to it.
var Unreachable *Terminator

// AppendSuccessors appends to the given slice any successors for the recieving
// terminator. Pass a nil slice to force this function to allocate a new backing
// array and return it, or pre-allocate a buffer in the caller.
//
// Most terminators have no more than two successors, so passing a slice with
// at least capacity two can avoid allocation in many cases. On the other hand,
// some terminators have no successors at all, so passing nil can mean avoiding
// allocation altogether in those cases.
func (t *Terminator) AppendSuccessors(to []*BasicBlock) []*BasicBlock {
	// This switch must cover all of the ops that are considered to be
	// terminator operations by op.Terminator.
	switch t.op {
	case OpJump:
		return append(to, t.args[0].Block)
	case OpBranch:
		return append(to, t.args[0].Block, t.args[1].Block)
	case OpSwitch:
		ret := to
		for _, arg := range t.args {
			ret = append(ret, arg.Block)
		}
		return ret
	case OpReturn, OpUnreachable:
		return to // no successors
	case OpYield, OpAwait:
		return append(to, t.args[0].Block)
	default:
		if t.op.Terminator() {
			// Indicates we're missing a case above
			panic(fmt.Sprintf("AppendSuccessors is missing a case for %s", t.op))
		} else {
			// Indicates an incorrectly-constructed terminator
			panic("AppendSuccessors with non-terminator operation")
		}
	}
}

// AddSuccessors adds to the given set any successors for the receiving
// terminator, in-place.
func (t *Terminator) AddSuccessors(to BasicBlockSet) {
	// For now we're going to implement this in terms of AppendSuccessors, which
	// requires us to allocate a backing array for this slice. We may wish to
	// rework this later to remove this allocation if it proves to be troublesome
	// after profiling.
	succs := t.AppendSuccessors(nil)
	for _, block := range succs {
		to.Add(block)
	}
}

// bufForArgs returns a zero-length arg slice with at least the given capacity
// that can be used as the arguments for the receiving terminator.
//
// bufForArgs either allocates a slice with the given capacity ready to have
// arguments appended to it, or returns a slice backed by t.argsBuf if
// its length is enough to contain the args with no further allocation.
func (t *Terminator) bufForArgs(capacity int) []BasicBlockValue {
	if len(t.argsBuf) >= capacity {
		return t.argsBuf[:0]
	}
	return make([]BasicBlockValue, 0, capacity)
}

func init() {
	Unreachable = &Terminator{
		op: OpUnreachable,
	}
}
