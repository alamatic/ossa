package oana

import (
	"testing"

	"github.com/alamatic/ossa"
)

func TestFindNaturalLoops(t *testing.T) {
	entry := &ossa.BasicBlock{}
	loopHeader := &ossa.BasicBlock{}
	loopBody := &ossa.BasicBlock{}
	loopTail := &ossa.BasicBlock{}
	exit := &ossa.BasicBlock{}

	entry.Terminator = ossa.Jump(loopHeader)
	loopHeader.Terminator = ossa.Branch(
		ossa.AuxLiteral(nil),
		loopBody,
		exit,
	)
	loopBody.Terminator = ossa.Jump(loopTail)
	loopTail.Terminator = ossa.Jump(loopHeader)
	exit.Terminator = ossa.Return(ossa.AuxLiteral(nil))

	preds := FindPredecessors(entry)
	doms := FindDominators(entry, preds)
	loops := FindNaturalLoops(doms, nil)

	// We care about the identities of these blocks rather than their contents,
	// so to make test results easier to understand we'll give each block a
	// name and compare by those names.
	names := map[*ossa.BasicBlock]string{
		entry:      "entry",
		loopHeader: "loopHeader",
		loopBody:   "loopBody",
		loopTail:   "loopTail",
		exit:       "exit",
	}

	got := loops
	want := []NaturalLoop{
		{Head: loopHeader, Tail: loopTail},
	}
	if len(got) != len(want) {
		t.Fatalf("wrong number of loops %d; want %d", len(got), len(want))
	}
	for i := range want {
		gotLoop := got[i]
		wantLoop := want[i]
		if gotLoop.Head != wantLoop.Head {
			t.Errorf("loop %d has wrong head %q; want %q", i, names[gotLoop.Head], names[wantLoop.Head])
		}
		if gotLoop.Tail != wantLoop.Tail {
			t.Errorf("loop %d has wrong tail %q; want %q", i, names[gotLoop.Tail], names[wantLoop.Tail])
		}
	}

	gotBody := loops[0].FindBody(preds)
	wantBody := ossa.NewBasicBlockSet(loopHeader, loopBody, loopTail)
	for b := range wantBody {
		if !gotBody.Has(b) {
			t.Errorf("loop body should contain %q", names[b])
		}
	}
	for b := range gotBody {
		if !wantBody.Has(b) {
			t.Errorf("loop body should not contain %q", names[b])
		}
	}
}
