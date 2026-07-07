package main

import "fmt"

// stackADT is the common interface every implementation satisfies, so main()
// can drive the same official operation sequence through each one.
type stackADT interface {
	Push(x int)
	Pop() int
	Top() int
	Empty() bool
}

// queue is a minimal FIFO queue backed by a slice, exposing only the
// operations a real queue would (enqueue at back, dequeue at front, peek,
// size). Both implementations below build a LIFO stack out of these.
type queue struct {
	data []int
}

func (q *queue) enqueue(x int) { q.data = append(q.data, x) } // add to back
func (q *queue) dequeue() int { // remove & return front
	x := q.data[0]
	q.data = q.data[1:]
	return x
}
func (q *queue) front() int { return q.data[0] } // peek front
func (q *queue) size() int  { return len(q.data) }

// ── Approach 1: Two Queues, Costly Push ──────────────────────────────────────
//
// TwoQueueStack implements a LIFO stack with two FIFO queues, making Push the
// expensive operation so that Pop/Top are O(1).
//
// Intuition:
//
//	A queue removes from the front, a stack from the back. To expose the
//	most-recent element at the front of q1, enqueue the new value into the
//	empty q2, then drain all of q1 behind it; swap the queue names. Now q1's
//	front is the newest element — exactly what Pop/Top must return.
//
// Algorithm:
//
//	Push:  enqueue x into q2; move every element of q1 to q2; swap q1,q2.
//	Pop:   dequeue q1's front (the newest element).
//	Top:   peek q1's front.
//	Empty: q1 has size 0.
//
// Time:  Push O(n), Pop O(1), Top O(1), Empty O(1).
// Space: O(n) — the n elements, transiently split across two queues.
type TwoQueueStack struct {
	q1 *queue // holds elements with newest at the front
	q2 *queue // scratch queue used during Push
}

// NewTwoQueueStack builds an empty two-queue stack.
func NewTwoQueueStack() *TwoQueueStack {
	return &TwoQueueStack{q1: &queue{}, q2: &queue{}}
}

// Push inserts x so that it becomes q1's front (top of the stack).
func (s *TwoQueueStack) Push(x int) {
	s.q2.enqueue(x) // new element goes in first, so it ends up at the front
	for s.q1.size() > 0 {
		s.q2.enqueue(s.q1.dequeue()) // older elements queued behind it, order preserved
	}
	s.q1, s.q2 = s.q2, s.q1 // q2 now has the desired order; make it the main queue
}

// Pop removes and returns the top (q1's front).
func (s *TwoQueueStack) Pop() int { return s.q1.dequeue() }

// Top returns the top without removing it.
func (s *TwoQueueStack) Top() int { return s.q1.front() }

// Empty reports whether the stack has no elements.
func (s *TwoQueueStack) Empty() bool { return s.q1.size() == 0 }

// ── Approach 2: One Queue, Costly Push (Optimal / Idiomatic) ─────────────────
//
// OneQueueStack implements the stack with a single FIFO queue by rotating it
// after each Push so the newest element sits at the front.
//
// Intuition:
//
//	Enqueue x at the back, then rotate the queue by dequeuing-and-re-enqueuing
//	the (size−1) elements that were already there. Those older elements circle
//	around behind x, leaving x at the front — so the front always holds the
//	most-recent element and Pop/Top are simple front operations.
//
// Algorithm:
//
//	Push:  enqueue x; then dequeue+enqueue the previous (size−1) elements.
//	Pop:   dequeue the front.
//	Top:   peek the front.
//	Empty: size 0.
//
// Time:  Push O(n), Pop O(1), Top O(1), Empty O(1).
// Space: O(n) — a single queue holding all elements.
type OneQueueStack struct {
	q *queue // single queue; front is always the current stack top
}

// NewOneQueueStack builds an empty one-queue stack.
func NewOneQueueStack() *OneQueueStack {
	return &OneQueueStack{q: &queue{}}
}

// Push adds x and rotates so x becomes the front.
func (s *OneQueueStack) Push(x int) {
	s.q.enqueue(x)                      // x goes to the back...
	for i := 0; i < s.q.size()-1; i++ { // ...rotate the OLD elements around behind it
		s.q.enqueue(s.q.dequeue())
	}
}

// Pop removes and returns the top (front).
func (s *OneQueueStack) Pop() int { return s.q.dequeue() }

// Top returns the top (front) without removing it.
func (s *OneQueueStack) Top() int { return s.q.front() }

// Empty reports whether the stack is empty.
func (s *OneQueueStack) Empty() bool { return s.q.size() == 0 }

// runExample drives the single official example through one implementation and
// returns the output list in LeetCode's null/value format.
//
// Ops:  ["MyStack","push","push","top","pop","empty"]
// Args: [[],[1],[2],[],[],[]]
func runExample(newStack func() stackADT) string {
	st := newStack()
	out := []string{"null"} // constructor
	st.Push(1)              // push(1) → null
	out = append(out, "null")
	st.Push(2) // push(2) → null
	out = append(out, "null")
	out = append(out, fmt.Sprintf("%d", st.Top()))   // top()  → 2
	out = append(out, fmt.Sprintf("%d", st.Pop()))   // pop()  → 2
	out = append(out, fmt.Sprintf("%t", st.Empty())) // empty() → false
	return "[" + join(out, ", ") + "]"
}

// join concatenates parts with sep (small helper to avoid importing strings).
func join(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for _, p := range parts[1:] {
		out += sep + p
	}
	return out
}

func main() {
	fmt.Println("=== Approach 1: Two Queues, Costly Push ===")
	fmt.Println(runExample(func() stackADT { return NewTwoQueueStack() })) // [null, null, null, 2, 2, false]

	fmt.Println("=== Approach 2: One Queue, Costly Push (Optimal) ===")
	fmt.Println(runExample(func() stackADT { return NewOneQueueStack() })) // [null, null, null, 2, 2, false]
}
