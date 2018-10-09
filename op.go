package ossa

type Op int

const (
	opInvalid Op = iota

	OpGlobalSym
	OpLocalSym
	OpArgument
	OpAuxLiteral
	OpPhi

	OpLoad
	OpStore

	OpCall

	// we also have some internal-only operations used to deal with CFG-related
	// concerns. These are not visible to callers.
	opBasicBlock

	// This special value represents the split between value operations and
	// terminator operations.
	opEndValues

	OpJump
	OpBranch
	OpSwitch
	OpReturn
	OpYield
	OpAwait
	OpUnreachable

	opEndTerminators
)

//go:generate stringer -type Op

// Valid returns true if the receiving op is valid, which is to say it is one
// of the constant values defined in this package. The zero value of
// Op is not valid.
func (o Op) Valid() bool {
	// This excludes opInvalid, opEndValues and opEndTerminators, along with
	// any values greater than opEndTerminators that are not defined yet.
	return o.Value() || o.Terminator()
}

// Value returns true if the receiving op belongs to the set of operations
// used with Values, as opposed to Terminators.
func (o Op) Value() bool {
	return o > opInvalid && o < opEndValues
}

// Terminator returns true if the receiving op belongs to the set of operations
// used with Terminators, as opposed to Values.
func (o Op) Terminator() bool {
	return o > opEndValues && o < opEndTerminators
}

// assertValue panics if the reciever is not a value
func (o Op) assertValue() {
	if !o.Value() {
		panic("operation is not suitable for value")
	}
}

// assertValue panics if the reciever is not a value
func (o Op) assertTerminator() {
	if !o.Terminator() {
		panic("operation is not suitable for terminator")
	}
}
