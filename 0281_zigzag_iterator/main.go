package main

import "fmt"

// zigzagADT is the common interface every approach implements so main() can
// drive the same official example through all of them.
type zigzagADT interface {
	HasNext() bool
	Next() int
}

// ── Approach 1: Flatten Upfront (Round-Robin Merge) ──────────────────────────
//
// FlattenZigzag solves Zigzag Iterator by pre-merging both lists into a single
// slice in zigzag order at construction time, then iterating that slice.
//
// Intuition:
//
//	The output is just the two lists interleaved: take one from list 1, one
//	from list 2, repeat; when one runs dry, keep draining the other. If we
//	build that merged order once up front, Next/HasNext become a trivial slice
//	walk. Simple and correct — the price is O(total) extra memory and no lazy
//	evaluation.
//
// Algorithm:
//
//	Constructor: walk index j = 0, 1, 2, …; at each j append v1[j] if it
//	             exists then v2[j] if it exists, until both are exhausted.
//	HasNext:     pos < len(merged).
//	Next:        return merged[pos] and advance pos.
//
// Time:  Constructor O(n+m); Next O(1); HasNext O(1).
// Space: O(n+m) — the flattened buffer.
type FlattenZigzag struct {
	merged []int // both lists interleaved in zigzag order
	pos    int   // index of the next element to emit
}

// NewFlattenZigzag builds the flattened structure from the two input lists.
func NewFlattenZigzag(v1, v2 []int) *FlattenZigzag {
	merged := []int{}
	for j := 0; j < len(v1) || j < len(v2); j++ {
		if j < len(v1) {
			merged = append(merged, v1[j]) // column j of list 1
		}
		if j < len(v2) {
			merged = append(merged, v2[j]) // column j of list 2
		}
	}
	return &FlattenZigzag{merged: merged}
}

// HasNext reports whether an unread element remains.
func (z *FlattenZigzag) HasNext() bool { return z.pos < len(z.merged) }

// Next returns the next element in zigzag order.
func (z *FlattenZigzag) Next() int {
	v := z.merged[z.pos] // element at the current cursor
	z.pos++              // advance
	return v
}

// ── Approach 2: Two Cursors, Turn Toggle ─────────────────────────────────────
//
// TwoCursorZigzag solves Zigzag Iterator lazily by keeping an index into each
// list and a "whose turn" flag, advancing only on demand.
//
// Intuition:
//
//	No need to buffer everything: keep a pointer per list and alternate turns.
//	Before emitting, skip the current list if it is exhausted so the turn
//	always lands on a list that still has data. This is O(1) extra space and
//	fully lazy — you can stop iterating early for free.
//
// Algorithm:
//
//	State: i1, i2 (per-list cursors), turn ∈ {0,1}.
//	HasNext: i1 < len(v1) || i2 < len(v2).
//	Next:    if turn == 0 and v1 still has data, take v1[i1++]; else take
//	         v2[i2++]. Flip turn (but if the chosen-next list is empty the
//	         skip logic keeps us on the non-empty one).
//
// Time:  Next O(1) amortised; HasNext O(1).
// Space: O(1) — only cursors and a flag.
type TwoCursorZigzag struct {
	v1, v2 []int // the two source lists
	i1, i2 int   // per-list read cursors
	turn   int   // 0 → prefer v1 next, 1 → prefer v2 next
}

// NewTwoCursorZigzag builds the lazy two-cursor structure.
func NewTwoCursorZigzag(v1, v2 []int) *TwoCursorZigzag {
	return &TwoCursorZigzag{v1: v1, v2: v2}
}

// HasNext reports whether either list still has unread elements.
func (z *TwoCursorZigzag) HasNext() bool {
	return z.i1 < len(z.v1) || z.i2 < len(z.v2)
}

// Next returns the next element, honouring turns but skipping exhausted lists.
func (z *TwoCursorZigzag) Next() int {
	// Take from v1 when it's v1's turn AND v1 has data, OR when v2 is empty.
	if (z.turn == 0 && z.i1 < len(z.v1)) || z.i2 >= len(z.v2) {
		v := z.v1[z.i1]
		z.i1++
		z.turn = 1 // next time prefer v2
		return v
	}
	v := z.v2[z.i2]
	z.i2++
	z.turn = 0 // next time prefer v1
	return v
}

// ── Approach 3: Queue of Iterators (k-List Generalisation, Optimal) ──────────
//
// QueueZigzag solves Zigzag Iterator with a round-robin queue of live cursors,
// which extends unchanged to k lists (the classic follow-up).
//
// Intuition:
//
//	Model each list as a cursor (slice + index). Keep a FIFO queue of cursors
//	that still have elements. To emit: pop the front cursor, take its current
//	value, advance it, and push it back only if it still has data. The queue
//	naturally cycles list 1 → list 2 → … → list k → list 1, i.e. zigzag, and
//	drops lists as they empty. Works for any number of lists.
//
// Algorithm:
//
//	Constructor: enqueue every non-empty list's cursor.
//	HasNext:     queue is non-empty.
//	Next:        dequeue cursor c; value = c.list[c.idx]; c.idx++; if c still
//	             has data, enqueue it again; return value.
//
// Time:  Next O(1); HasNext O(1).
// Space: O(k) cursors in the queue (k = number of lists).
type cursor struct {
	list []int // the list this cursor walks
	idx  int   // next index to read within list
}

type QueueZigzag struct {
	queue []*cursor // live cursors in round-robin order
}

// NewQueueZigzag builds the queue from any number of lists (here two).
func NewQueueZigzag(lists ...[]int) *QueueZigzag {
	z := &QueueZigzag{}
	for _, l := range lists {
		if len(l) > 0 {
			z.queue = append(z.queue, &cursor{list: l}) // only track non-empty lists
		}
	}
	return z
}

// HasNext reports whether any cursor still has data.
func (z *QueueZigzag) HasNext() bool { return len(z.queue) > 0 }

// Next emits the front cursor's value and re-queues it if not exhausted.
func (z *QueueZigzag) Next() int {
	c := z.queue[0]       // front cursor (whose turn it is)
	z.queue = z.queue[1:] // dequeue it
	v := c.list[c.idx]    // current value
	c.idx++               // advance this cursor
	if c.idx < len(c.list) {
		z.queue = append(z.queue, c) // still has data → back of the line
	}
	return v
}

// drain runs the iterator to completion and returns the emitted order.
func drain(z zigzagADT) []int {
	out := []int{}
	for z.HasNext() {
		out = append(out, z.Next())
	}
	return out
}

func main() {
	v1 := []int{1, 2}
	v2 := []int{3, 4, 5, 6}

	fmt.Println("=== Approach 1: Flatten Upfront ===")
	fmt.Println(drain(NewFlattenZigzag(v1, v2))) // [1 3 2 4 5 6]

	fmt.Println("=== Approach 2: Two Cursors, Turn Toggle ===")
	fmt.Println(drain(NewTwoCursorZigzag(v1, v2))) // [1 3 2 4 5 6]

	fmt.Println("=== Approach 3: Queue of Iterators (Optimal) ===")
	fmt.Println(drain(NewQueueZigzag(v1, v2))) // [1 3 2 4 5 6]
}
