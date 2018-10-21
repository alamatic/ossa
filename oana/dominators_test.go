package oana

import (
	"testing"

	"github.com/alamatic/ossa"
)

func TestFindDominators(t *testing.T) {
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

	preds := FindPredecessors(entry)
	doms := FindDominators(entry, preds)

	// We care about the identities of these blocks rather than their contents,
	// so to make test results easier to understand we'll give each block a
	// name and compare by those names.
	names := map[*ossa.BasicBlock]string{
		entry:      "entry",
		loopHeader: "loopHeader",
		loopBody:   "loopBody",
		exit:       "exit",
	}

	got := doms
	want := DominatorsTable{
		entry:      ossa.NewBasicBlockSet(entry),
		loopHeader: ossa.NewBasicBlockSet(entry, loopHeader),
		loopBody:   ossa.NewBasicBlockSet(entry, loopHeader, loopBody),
		exit:       ossa.NewBasicBlockSet(entry, loopHeader, exit),
	}
	for wantB, wantDBs := range want {
		gotDBs := got[wantB]
		for wantDB := range wantDBs {
			if !gotDBs.Has(wantDB) {
				t.Errorf("%q should dominate %q", names[wantDB], names[wantB])
			}
		}
		for gotDB := range gotDBs {
			if !wantDBs.Has(gotDB) {
				t.Errorf("%q should not dominate %q", names[gotDB], names[wantB])
			}
		}
	}
	for gotB := range got {
		if _, exists := want[gotB]; !exists {
			t.Errorf("%q should not be in the result", names[gotB])
		}
	}
}
