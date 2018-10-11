package oana

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/alamatic/ossa"
)

func TestForwardDataFlow(t *testing.T) {
	entry := &ossa.BasicBlock{}
	loopHeader := &ossa.BasicBlock{}
	loopBody := &ossa.BasicBlock{}
	exit := &ossa.BasicBlock{}

	entry.Terminator = ossa.Jump(loopHeader)
	loopHeader.Terminator = ossa.Branch(
		ossa.AuxLiteral(nil),
		loopBody,
		exit,
	)
	loopBody.Terminator = ossa.Jump(loopHeader)
	exit.Terminator = ossa.Return(ossa.AuxLiteral(nil))

	a := &loggingBlockAnalyzer{
		// We'll simulate a typical situation where the first visit to each
		// block causes a change and then, when we revisit the loopHeader,
		// it has a further change caused by the additional information from
		// visiting the loop body, after which everything reaches fixpoint.
		changeCount: map[*ossa.BasicBlock]int{
			entry:      1,
			loopHeader: 2,
			loopBody:   1,
			exit:       1,
		},
	}

	ForwardDataFlow(entry, a)

	// We care about the identities of these blocks rather than their contents,
	// so to make test results easier to understand we'll give each block a
	// name and compare by those names.
	names := map[*ossa.BasicBlock]string{
		entry:      "entry",
		loopHeader: "loopHeader",
		loopBody:   "loopBody",
		exit:       "exit",
	}

	got := make([]string, len(a.calls))
	for i, block := range a.calls {
		got[i] = names[block]
	}
	want := []string{
		"entry",
		"loopHeader",
		"loopBody",
		"loopHeader", // visited again because loopBody points back at it
		"loopBody",   // visited one more time after loopHeader, but finds fixpoint
		"exit",       // we reach exit only after the loop blocks have reached fixpoint
	}
	if !cmp.Equal(got, want) {
		t.Errorf("wrong block visit order\ngot: %#v\nwant: %#v", got, want)
	}
}

type loggingBlockAnalyzer struct {
	changeCount map[*ossa.BasicBlock]int
	calls       []*ossa.BasicBlock
}

func (a *loggingBlockAnalyzer) AnalyzeBlock(block *ossa.BasicBlock) (changed bool) {
	a.calls = append(a.calls, block)
	if ct := a.changeCount[block]; ct > 0 {
		a.changeCount[block] = ct - 1
		return true
	}
	return false
}
