package oana

import (
	"github.com/alamatic/ossa"
)

type blockQueue interface {
	Add(block *ossa.BasicBlock)
	Has(block *ossa.BasicBlock) bool
	Empty() bool
	Peek() *ossa.BasicBlock
	Next() *ossa.BasicBlock
}

// blockLIFO is a stack data structure that preserves insertion order while
// also guaranteeing that the same item cannot appear twice in the stack.
//
// This data structure is not safe for concurrent modifications or reads
// concurrent with modifications.
type blockLIFO struct {
	items   []*ossa.BasicBlock
	present ossa.BasicBlockSet
}

var _ blockQueue = (*blockLIFO)(nil)

// newBlockLIFO allocates a new LIFO stack with the given initial capacity.
// If the length grows beyond this initial capacity then a new buffer will
// be allocated, growing the capacity.
func newBlockLIFO(initialCapacity int) *blockLIFO {
	var items []*ossa.BasicBlock
	items = make([]*ossa.BasicBlock, 0, initialCapacity)
	return &blockLIFO{
		items:   items,
		present: make(ossa.BasicBlockSet),
	}
}

// Add ensures that the given block is present in the stack. If it is already
// present, no action is taken. If it is not already present then it is
// pushed on the top of the stack.
//
// This is an implementation of ossa.BasicBlockAdder, so a block stack can be
// used with functions that can add blocks to a collection via this interface.
func (q *blockLIFO) Add(block *ossa.BasicBlock) {
	if q.present.Has(block) {
		return // already in the queue
	}
	q.items = append(q.items, block)
	q.present.Add(block)
}

// Has tests whether the given block is already present in the stack, returning
// true if so.
func (q *blockLIFO) Has(block *ossa.BasicBlock) bool {
	return q.present.Has(block)
}

// Empty returns true if stack is empty, and false otherwise.
func (q *blockLIFO) Empty() bool {
	return len(q.items) == 0
}

// Length returns te number of items in the stack.
func (q *blockLIFO) Length() int {
	return len(q.items)
}

// ReverseTopN reverses the order of the top n items in the stack, in-place. Will
// panic if there aren't at least that many items to reverse.
//
// This is a rather specialized method intended to help algorithms that
// wish to visit block successors in the opposite order that they are generated
// by the ossa.Terminator type, in conjunction with Length to find how many
// new items were added.
func (q *blockLIFO) ReverseTopN(n int) {
	l := len(q.items)
	items := q.items[l-n:]
	l = len(items)
	c := l / 2
	for i := 0; i < c; i++ {
		items[i], items[l-i-1] = items[l-i-1], items[i]
	}
}

// Peek returns the next item in the stack without taking it, or returns nil
// if the stack is currently empty.
func (q *blockLIFO) Peek() *ossa.BasicBlock {
	if q.Empty() {
		return nil
	}
	return q.items[len(q.items)-1]
}

// Next removes the top item from the stack and returns it. It returns nil
// if the stack is currently empty.
func (q *blockLIFO) Next() *ossa.BasicBlock {
	ret := q.Peek()
	if ret == nil {
		return nil
	}
	q.items = q.items[:len(q.items)-1]
	q.present.Remove(ret)
	return ret
}

// blockFIFO is a queue data structure that preserves insertion order while
// also guaranteeing that the same item cannot appear twice in the queue.
//
// This data structure is not safe for concurrent modifications or reads
// concurrent with modifications.
type blockFIFO struct {
	next, end *blockFIFOEntry
	present   ossa.BasicBlockSet
}

var _ blockQueue = (*blockFIFO)(nil)

type blockFIFOEntry struct {
	block *ossa.BasicBlock
	next  *blockFIFOEntry
}

func newblockFIFO() *blockFIFO {
	return &blockFIFO{
		present: make(ossa.BasicBlockSet),
	}
}

// Add ensures that the given block is present in the queue. If it is already
// present, no action is taken. If it is not already present then it is
// appended to the end of the queue.
//
// This is an implementation of ossa.BasicBlockAdder, so a block queue can be
// used with functions that can add blocks to a collection via this interface.
func (q *blockFIFO) Add(block *ossa.BasicBlock) {
	if q.present.Has(block) {
		return // already in the queue
	}

	var entry blockFIFOEntry
	entry.block = block
	if q.end != nil {
		q.end.next = &entry
	} else {
		// Adding to an empty queue
		q.next = &entry
	}
	q.end = &entry
	q.present.Add(block)
}

// Has tests whether the given block is already present in the queue, returning
// true if so.
func (q *blockFIFO) Has(block *ossa.BasicBlock) bool {
	return q.present.Has(block)
}

// Empty returns true if queue is empty, and false otherwise.
func (q *blockFIFO) Empty() bool {
	return q.next == nil
}

// Peek returns the next item in the queue without taking it, or returns nil
// if the queue is currently empty.
func (q *blockFIFO) Peek() *ossa.BasicBlock {
	return q.next.block
}

// Next removes the next item from the queue and returns it. It returns nil
// if the queue is currently empty.
func (q *blockFIFO) Next() *ossa.BasicBlock {
	if q.next == nil {
		return nil
	}

	ret := q.next.block
	q.next = q.next.next
	if q.next == nil {
		// Queue is now empty, so it has no end either
		q.end = nil
	}
	q.present.Remove(ret)
	return ret
}
