package main

import "fmt"

// queueADT is the common interface each implementation satisfies so main() can
// drive the same official operation sequence through every approach.
type queueADT interface {
	Push(x int)  // enqueue at the back
	Pop() int    // dequeue from the front, returning the removed element
	Peek() int   // return the front element without removing it
	Empty() bool // report whether the queue holds no elements
}

// ── Approach 1: Two Stacks, Costly Push ──────────────────────────────────────
//
// CostlyPushQueue implements a FIFO queue with two LIFO stacks, keeping the
// "in" stack always ordered so the front is on top.
//
// Intuition:
//
//	A stack reverses insertion order; queueing needs to preserve it. If on
//	every push we move everything to a helper, drop the new element at the
//	bottom, then move it all back, the main stack always has the oldest
//	element on top — so Pop/Peek are O(1). We pay the reordering cost up
//	front, at push time.
//
// Algorithm:
//
//	Push: pop all of `in` into `out`; push x onto `in`; pop all of `out`
//	      back onto `in`. Now `in` is front-on-top.
//	Pop/Peek: read/remove the top of `in`.
//
// Time:  Push O(n), Pop O(1), Peek O(1), Empty O(1).
// Space: O(n) — the two stacks together hold n elements.
type CostlyPushQueue struct {
	in  []int // primary stack, maintained with the queue front on top
	out []int // scratch stack used only during Push
}

// NewCostlyPushQueue builds an empty queue.
func NewCostlyPushQueue() *CostlyPushQueue {
	return &CostlyPushQueue{}
}

// Push inserts x at the back while keeping the front on top of `in`.
func (q *CostlyPushQueue) Push(x int) {
	for len(q.in) > 0 { // dump existing elements onto the scratch stack
		q.out = append(q.out, q.in[len(q.in)-1])
		q.in = q.in[:len(q.in)-1]
	}
	q.in = append(q.in, x) // new element goes in first → ends up at the bottom
	for len(q.out) > 0 {   // restore the old elements on top of it
		q.in = append(q.in, q.out[len(q.out)-1])
		q.out = q.out[:len(q.out)-1]
	}
}

// Pop removes and returns the front element (top of `in`).
func (q *CostlyPushQueue) Pop() int {
	top := q.in[len(q.in)-1] // front element sits on top by invariant
	q.in = q.in[:len(q.in)-1]
	return top
}

// Peek returns the front element without removing it.
func (q *CostlyPushQueue) Peek() int {
	return q.in[len(q.in)-1] // front is always the top of `in`
}

// Empty reports whether the queue has no elements.
func (q *CostlyPushQueue) Empty() bool {
	return len(q.in) == 0
}

// ── Approach 2: Two Stacks, Amortized O(1) Push (Optimal) ────────────────────
//
// AmortizedQueue implements a FIFO queue with an input stack and an output
// stack, transferring lazily so each element moves at most once.
//
// Intuition:
//
//	Keep pushes cheap: just drop new elements on an `in` stack. When we need
//	the front and the `out` stack is empty, pour all of `in` into `out`.
//	Reversing a LIFO stack into another LIFO stack yields FIFO order, so the
//	top of `out` is the oldest element. Each element is moved from `in` to
//	`out` exactly once over its lifetime → amortized O(1) per operation.
//
// Algorithm:
//
//	Push: append x to `in`.
//	transfer(): if `out` is empty, pop all of `in` onto `out`.
//	Pop:  transfer(); pop the top of `out`.
//	Peek: transfer(); read the top of `out`.
//	Empty: both stacks empty.
//
// Time:  Push O(1), Pop O(1) amortized, Peek O(1) amortized, Empty O(1).
// Space: O(n).
type AmortizedQueue struct {
	in  []int // newest elements, back of the queue on top
	out []int // oldest elements reversed, front of the queue on top
}

// NewAmortizedQueue builds an empty queue.
func NewAmortizedQueue() *AmortizedQueue {
	return &AmortizedQueue{}
}

// Push simply appends to the input stack (deferred ordering).
func (q *AmortizedQueue) Push(x int) {
	q.in = append(q.in, x) // O(1): reordering happens later, on demand
}

// transfer refills `out` from `in` only when `out` has run dry.
func (q *AmortizedQueue) transfer() {
	if len(q.out) == 0 { // only reverse when the front supply is exhausted
		for len(q.in) > 0 {
			q.out = append(q.out, q.in[len(q.in)-1]) // reverse in → out = FIFO
			q.in = q.in[:len(q.in)-1]
		}
	}
}

// Pop moves the front to `out` if needed, then removes and returns it.
func (q *AmortizedQueue) Pop() int {
	q.transfer()               // ensure the front is on top of `out`
	top := q.out[len(q.out)-1] // oldest element
	q.out = q.out[:len(q.out)-1]
	return top
}

// Peek returns the front element without removing it.
func (q *AmortizedQueue) Peek() int {
	q.transfer() // make sure `out` holds the front
	return q.out[len(q.out)-1]
}

// Empty reports whether both stacks are empty.
func (q *AmortizedQueue) Empty() bool {
	return len(q.in) == 0 && len(q.out) == 0 // no element anywhere
}

// runExample drives the single official LeetCode example through one queue
// implementation and returns the output list in LeetCode's format.
//
// Ops:  ["MyQueue","push","push","peek","pop","empty"]
// Args: [[],[1],[2],[],[],[]]
func runExample(newQueue func() queueADT) string {
	q := newQueue()         // "MyQueue" → null
	out := []string{"null"} // constructor produces no value
	q.Push(1)               // push(1) → null
	out = append(out, "null")
	q.Push(2) // push(2) → null
	out = append(out, "null")
	out = append(out, fmt.Sprintf("%d", q.Peek()))  // peek() → 1
	out = append(out, fmt.Sprintf("%d", q.Pop()))   // pop()  → 1
	out = append(out, fmt.Sprintf("%t", q.Empty())) // empty()→ false

	res := "["
	for i, v := range out {
		if i > 0 {
			res += ", "
		}
		res += v
	}
	return res + "]"
}

func main() {
	fmt.Println("=== Approach 1: Two Stacks, Costly Push ===")
	fmt.Println(runExample(func() queueADT { return NewCostlyPushQueue() })) // [null, null, null, 1, 1, false]

	fmt.Println("=== Approach 2: Two Stacks, Amortized O(1) Push (Optimal) ===")
	fmt.Println(runExample(func() queueADT { return NewAmortizedQueue() })) // [null, null, null, 1, 1, false]
}
