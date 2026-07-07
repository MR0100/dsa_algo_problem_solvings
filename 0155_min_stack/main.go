package main

import (
	"fmt"
	"strings"
)

// minStacker lets main() drive every implementation through one test harness.
type minStacker interface {
	Push(val int)
	Pop()
	Top() int
	GetMin() int
}

// ── Approach 1: Brute Force (Single Stack, Scan for Min) ────────────────────
//
// BruteForceMinStack solves Min Stack with one plain slice; GetMin rescans
// the whole stack on every call.
//
// Intuition:
//
//	Push/Pop/Top are what a slice already gives us for free. The only hard
//	operation is GetMin — the lazy answer is to recompute it on demand by
//	scanning all live elements. Correct, but violates the O(1)-per-operation
//	requirement.
//
// Algorithm:
//
//	Push: append val.  Pop: drop the last element.  Top: return the last
//	element.  GetMin: linear scan over the slice keeping the smallest.
//
// Time:  Push/Pop/Top O(1); GetMin O(n).
// Space: O(n) — one slot per live element.
type BruteForceMinStack struct {
	data []int // plain stack storage; top is the last element
}

func NewBruteForceMinStack() *BruteForceMinStack { return &BruteForceMinStack{} }

func (s *BruteForceMinStack) Push(val int) { s.data = append(s.data, val) }

func (s *BruteForceMinStack) Pop() { s.data = s.data[:len(s.data)-1] } // drop top

func (s *BruteForceMinStack) Top() int { return s.data[len(s.data)-1] } // peek top

func (s *BruteForceMinStack) GetMin() int {
	m := s.data[0] // guaranteed non-empty by the problem statement
	for _, v := range s.data[1:] {
		if v < m {
			m = v // smaller element found during the rescan
		}
	}
	return m
}

// ── Approach 2: Stack of (Value, MinSoFar) Pairs ─────────────────────────────
//
// PairMinStack solves Min Stack by storing, with every element, the minimum
// of the stack at the moment that element was pushed.
//
// Intuition:
//
//	The minimum of a stack only changes when elements are pushed or popped —
//	and popping simply RESTORES whatever the minimum was before the push.
//	So snapshot "min of everything below and including me" into each entry;
//	the top entry's snapshot is always the current minimum, and Pop rolls the
//	state back automatically by discarding the top snapshot.
//
// Algorithm:
//
//	Push: min = val if stack empty, else min(val, top.min); push {val, min}.
//	Pop: drop the top pair.  Top: top pair's val.  GetMin: top pair's min.
//
// Time:  O(1) for every operation — no scanning, ever.
// Space: O(n) — two ints per element (2n words).
type PairMinStack struct {
	data []pair // each entry remembers the min at its push time
}

// pair bundles a stored value with the stack minimum as of its push.
type pair struct {
	val int // the element itself
	min int // min of the stack up to and including this element
}

func NewPairMinStack() *PairMinStack { return &PairMinStack{} }

func (s *PairMinStack) Push(val int) {
	m := val // an empty stack's new minimum is the pushed value itself
	if len(s.data) > 0 && s.data[len(s.data)-1].min < val {
		m = s.data[len(s.data)-1].min // previous min still smaller → keep it
	}
	s.data = append(s.data, pair{val: val, min: m})
}

func (s *PairMinStack) Pop() { s.data = s.data[:len(s.data)-1] } // snapshot rolls back

func (s *PairMinStack) Top() int { return s.data[len(s.data)-1].val }

func (s *PairMinStack) GetMin() int { return s.data[len(s.data)-1].min }

// ── Approach 3: Two Stacks (Lazy Min Stack) ──────────────────────────────────
//
// TwoStacksMinStack solves Min Stack with a main stack plus an auxiliary
// stack that only records elements which are new minima (<= current min).
//
// Intuition:
//
//	Approach 2 wastes space repeating the same min for long stretches.
//	Observation: the sequence of stack minima over time is exactly the
//	sequence of pushed values that were <= the min at their push time.
//	Keep just those in a second stack; its top IS the current minimum.
//	On Pop, if the popped value equals the top of the min stack, that
//	minimum's reign ends — pop it too, revealing the previous minimum.
//
// Algorithm:
//
//	Push: append to data; if mins empty or val <= mins top, also append to
//	      mins (the <= is vital for duplicate minima).
//	Pop:  if popped value == mins top, pop mins as well.
//	Top:  data top.  GetMin: mins top.
//
// Time:  O(1) for every operation.
// Space: O(n) worst case (strictly decreasing pushes), but only one extra slot per RECORD-LOW element in practice.
type TwoStacksMinStack struct {
	data []int // all elements
	mins []int // history of minima; top = current minimum
}

func NewTwoStacksMinStack() *TwoStacksMinStack { return &TwoStacksMinStack{} }

func (s *TwoStacksMinStack) Push(val int) {
	s.data = append(s.data, val)
	// push onto mins only when val ties or beats the current minimum;
	// the <= (not <) makes duplicate minima pop back correctly one by one
	if len(s.mins) == 0 || val <= s.mins[len(s.mins)-1] {
		s.mins = append(s.mins, val)
	}
}

func (s *TwoStacksMinStack) Pop() {
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1] // remove from the main stack
	if v == s.mins[len(s.mins)-1] {
		s.mins = s.mins[:len(s.mins)-1] // its reign as minimum is over
	}
}

func (s *TwoStacksMinStack) Top() int { return s.data[len(s.data)-1] }

func (s *TwoStacksMinStack) GetMin() int { return s.mins[len(s.mins)-1] }

// ── Approach 4: Difference Encoding (Optimal — One Stack, O(1) Extra) ───────
//
// DiffMinStack solves Min Stack with a single stack of DIFFERENCES from the
// running minimum plus one scalar holding the current minimum.
//
// Intuition:
//
//	Storing (value, min) pairs doubles memory. Trick: store val - min
//	instead of val. Then the sign of the stored difference encodes history:
//	  diff >= 0 → val did not change the min; val = min + diff.
//	  diff <  0 → val became the NEW min (so val == current min), and the
//	              PREVIOUS min can be recovered as min - diff.
//	One scalar (min) plus one number per element reconstructs everything —
//	the previous minima are woven into the stack itself.
//	Differences are stored as int64 because val - min can reach 2^32 - 1,
//	overflowing 32-bit ints (fine for Go's 64-bit int, but int64 makes the
//	requirement explicit and portable).
//
// Algorithm:
//
//	Push: empty → store 0, min = val.
//	      else  → store val - min; if val < min, min = val.
//	Pop:  d = pop; if d < 0 the popped element was the min → restore
//	      min = min - d (the old min).
//	Top:  d = peek; d < 0 → top IS the min → return min; else min + d.
//	GetMin: return min.
//
// Time:  O(1) for every operation.
// Space: O(n) for the stack itself, O(1) EXTRA — one int64 per element plus a single scalar, versus 2 per element in Approach 2.
type DiffMinStack struct {
	diffs []int64 // stored as val - minAtPushTime (int64: avoids overflow)
	min   int64   // current minimum of the whole stack
}

func NewDiffMinStack() *DiffMinStack { return &DiffMinStack{} }

func (s *DiffMinStack) Push(val int) {
	v := int64(val)
	if len(s.diffs) == 0 {
		s.diffs = append(s.diffs, 0) // first element: diff 0, min = itself
		s.min = v
		return
	}
	s.diffs = append(s.diffs, v-s.min) // negative ⇔ new record low
	if v < s.min {
		s.min = v // val is the new minimum
	}
}

func (s *DiffMinStack) Pop() {
	d := s.diffs[len(s.diffs)-1]
	s.diffs = s.diffs[:len(s.diffs)-1]
	if d < 0 {
		// popped element was the current min (stored while min changed);
		// previous min = current min - d  (d negative → subtracting grows it)
		s.min = s.min - d
	}
}

func (s *DiffMinStack) Top() int {
	d := s.diffs[len(s.diffs)-1]
	if d < 0 {
		return int(s.min) // negative diff ⇒ this element IS the current min
	}
	return int(s.min + d) // reconstruct: val = min + diff
}

func (s *DiffMinStack) GetMin() int { return int(s.min) }

// runExample drives one implementation through the official LeetCode example:
//
//	ops:  ["MinStack","push","push","push","getMin","pop","top","getMin"]
//	args: [[],[-2],[0],[-3],[],[],[],[]]
//
// and prints the produced output list, expected: [null,null,null,null,-3,null,0,-2]
func runExample(s minStacker) {
	out := []string{"null"} // constructor produces null
	s.Push(-2)
	out = append(out, "null")
	s.Push(0)
	out = append(out, "null")
	s.Push(-3)
	out = append(out, "null")
	out = append(out, fmt.Sprint(s.GetMin())) // -3
	s.Pop()
	out = append(out, "null")
	out = append(out, fmt.Sprint(s.Top()))    // 0
	out = append(out, fmt.Sprint(s.GetMin())) // -2
	fmt.Printf("[%s]\n", strings.Join(out, ","))
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Single Stack, Scan for Min) ===")
	runExample(NewBruteForceMinStack()) // [null,null,null,null,-3,null,0,-2]

	fmt.Println("=== Approach 2: Stack of (Value, MinSoFar) Pairs ===")
	runExample(NewPairMinStack()) // [null,null,null,null,-3,null,0,-2]

	fmt.Println("=== Approach 3: Two Stacks (Lazy Min Stack) ===")
	runExample(NewTwoStacksMinStack()) // [null,null,null,null,-3,null,0,-2]

	fmt.Println("=== Approach 4: Difference Encoding (Optimal) ===")
	runExample(NewDiffMinStack()) // [null,null,null,null,-3,null,0,-2]
}
