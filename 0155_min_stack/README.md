# 0155 — Min Stack

> LeetCode #155 · Difficulty: Medium
> **Categories:** Stack, Design

---

## Problem Statement

Design a stack that supports push, pop, top, and retrieving the minimum element in constant time.

Implement the `MinStack` class:

- `MinStack()` initializes the stack object.
- `void push(int val)` pushes the element `val` onto the stack.
- `void pop()` removes the element on the top of the stack.
- `int top()` gets the top element of the stack.
- `int getMin()` retrieves the minimum element in the stack.

You must implement a solution with `O(1)` time complexity for each function.

**Example 1:**
```
Input
["MinStack","push","push","push","getMin","pop","top","getMin"]
[[],[-2],[0],[-3],[],[],[],[]]

Output
[null,null,null,null,-3,null,0,-2]

Explanation
MinStack minStack = new MinStack();
minStack.push(-2);
minStack.push(0);
minStack.push(-3);
minStack.getMin(); // return -3
minStack.pop();
minStack.top();    // return 0
minStack.getMin(); // return -2
```

**Constraints:**
- `-2^31 <= val <= 2^31 - 1`
- Methods `pop`, `top` and `getMin` operations will always be called on **non-empty** stacks.
- At most `3 * 10^4` calls will be made to `push`, `pop`, `top`, and `getMin`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Bloomberg  | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Apple      | ★★★☆☆ Medium     | 2024          |
| Uber       | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack** — LIFO storage is the substrate of the whole design; the key insight is that popping *restores* past state, so per-element snapshots suffice → see [`/dsa/stack.md`](/dsa/stack.md)
- **Auxiliary/monotonic min stack** — the history of record-low pushes forms a non-increasing stack whose top is always the current minimum → see [`/dsa/monotonic_stack.md`](/dsa/monotonic_stack.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (scan on getMin) | push/pop/top O(1), getMin O(n) | O(n) | Never — fails the O(1) requirement; baseline only |
| 2 | Pair stack (value, minSoFar) | O(1) all ops | O(n): 2 words/element | Simplest correct O(1) design; the go-to interview answer |
| 3 | Two stacks (lazy min stack) | O(1) all ops | O(n) worst, ≤1 extra word per record low | When pushes rarely set new minima — less memory in practice |
| 4 | Difference encoding | O(1) all ops | O(n): 1 word/element + 1 scalar | The space-optimal follow-up; shows off invariant reasoning |

---

## Approach 1 — Brute Force (Single Stack, Scan for Min)

### Intuition
A plain slice already gives `push`, `pop`, and `top` in O(1). The only non-trivial operation is `getMin` — the lazy answer recomputes it on demand by scanning every live element. Correct, but `getMin` costs O(n), violating the required O(1)-per-operation bound; it exists to motivate the real designs.

### Algorithm
1. `Push(val)`: append `val` to the slice.
2. `Pop()`: truncate the last element.
3. `Top()`: return the last element.
4. `GetMin()`: scan the whole slice, tracking the smallest value; return it.

### Complexity
- **Time:** `push`/`pop`/`top` O(1) amortized; `getMin` O(n) — a full rescan per call.
- **Space:** O(n) — one word per live element.

### Code
```go
type BruteForceMinStack struct {
	data []int // plain stack storage; top is the last element
}

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
```

### Dry Run
Example 1 operation sequence:

| Op | data (bottom→top) | Returned |
|----|-------------------|----------|
| push(-2) | [-2] | null |
| push(0) | [-2, 0] | null |
| push(-3) | [-2, 0, -3] | null |
| getMin() | scan -2, 0, -3 | **-3** |
| pop() | [-2, 0] | null |
| top() | [-2, 0] | **0** |
| getMin() | scan -2, 0 | **-2** |

Output `[null,null,null,null,-3,null,0,-2]` ✓ (but each getMin walked the whole stack)

---

## Approach 2 — Stack of (Value, MinSoFar) Pairs

### Intuition
The stack minimum changes only on push and pop — and a pop simply *restores* whatever the minimum was before the corresponding push. That's undo semantics, and stacks are undo machines: store with every element a snapshot of "the minimum of everything up to and including me". The snapshot at the top is always the live minimum, and popping rolls state back for free by discarding the top snapshot.

### Algorithm
1. `Push(val)`: compute `m = val` if empty, else `m = min(val, top.min)`; push `{val, m}`.
2. `Pop()`: discard the top pair (the previous pair's snapshot automatically becomes current).
3. `Top()`: return top pair's `val`.
4. `GetMin()`: return top pair's `min`.

### Complexity
- **Time:** O(1) for every operation — each is a constant number of slice/struct accesses.
- **Space:** O(n) — two integers per element (2n words), regardless of data patterns.

### Code
```go
// pair bundles a stored value with the stack minimum as of its push.
type pair struct {
	val int // the element itself
	min int // min of the stack up to and including this element
}

type PairMinStack struct {
	data []pair // each entry remembers the min at its push time
}

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
```

### Dry Run
Example 1:

| Op | data as (val, min) pairs (bottom→top) | Returned |
|----|----------------------------------------|----------|
| push(-2) | [(-2,-2)] | null |
| push(0) | [(-2,-2), (0,-2)] | null |
| push(-3) | [(-2,-2), (0,-2), (-3,-3)] | null |
| getMin() | top pair (-3,**-3**) | **-3** |
| pop() | [(-2,-2), (0,-2)] | null |
| top() | top pair (**0**,-2) | **0** |
| getMin() | top pair (0,**-2**) | **-2** |

Output `[null,null,null,null,-3,null,0,-2]` ✓

---

## Approach 3 — Two Stacks (Lazy Min Stack)

### Intuition
Approach 2 repeats the same `min` snapshot for long stretches — wasteful. Notice *when* the minimum actually changes: exactly when a pushed value is `<=` the current minimum. Record only those record-low values in a second stack (`mins`); its top is always the current minimum, and it is naturally **non-increasing** from bottom to top (a monotonic stack). On `Pop`, if the value leaving equals the top of `mins`, that minimum's reign ends — pop `mins` too, exposing the previous record.

The `<=` (rather than `<`) when pushing onto `mins` is critical: with duplicate minima (push 1, push 1), each copy must have its own entry so that popping one duplicate doesn't wipe out the still-valid minimum.

### Algorithm
1. `Push(val)`: append to `data`; if `mins` is empty **or** `val <= mins.top`, append to `mins`.
2. `Pop()`: remove `v` from `data`; if `v == mins.top`, pop `mins` as well.
3. `Top()`: return `data` top.
4. `GetMin()`: return `mins` top.

### Complexity
- **Time:** O(1) per operation — at most one push/pop on each of the two stacks.
- **Space:** O(n) worst case (strictly decreasing input pushes every element onto `mins`), but only one extra word per *record-low* element in typical workloads — strictly no more than Approach 2 and usually far less.

### Code
```go
type TwoStacksMinStack struct {
	data []int // all elements
	mins []int // history of minima; top = current minimum
}

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
```

### Dry Run
Example 1:

| Op | data | mins | Returned |
|----|------|------|----------|
| push(-2) | [-2] | [-2] (empty → push) | null |
| push(0) | [-2, 0] | [-2] (0 > -2 → skip) | null |
| push(-3) | [-2, 0, -3] | [-2, -3] (-3 ≤ -2 → push) | null |
| getMin() | — | top = **-3** | **-3** |
| pop() | [-2, 0] (popped -3 == mins top → pop mins) | [-2] | null |
| top() | top = **0** | — | **0** |
| getMin() | — | top = **-2** | **-2** |

Output `[null,null,null,null,-3,null,0,-2]` ✓ — note `mins` never stored the 0.

---

## Approach 4 — Difference Encoding (Optimal — One Stack, O(1) Extra)

### Intuition
Can the old minima be recovered *without* storing them anywhere? Yes — weave them into the stack itself. Store `diff = val - min` (min taken *before* updating it) instead of `val`, and keep a single scalar `min`:

- `diff >= 0` → the push didn't change the minimum; the element's value is `min + diff`.
- `diff < 0` → the push **set a new record low**, so this element *is* the current minimum (`val == min`), and — the magic — the *previous* minimum is recoverable as `min - diff` (because `diff = val - oldMin` and `val = newMin`, so `oldMin = newMin - diff`).

Every pop of a negative diff therefore restores the exact previous minimum, all the way down. One number per element plus one scalar: strictly less state than any snapshot scheme.

Overflow note: `val - min` can reach `2^32 - 1` (e.g. `2^31-1 - (-2^31)`), which overflows a 32-bit int — the diffs are held as `int64` (Go's `int` is 64-bit on mainstream platforms, but `int64` makes the requirement explicit).

### Algorithm
1. `Push(val)`:
   - Empty stack → push `0`, set `min = val`.
   - Else push `val - min`; if `val < min`, set `min = val`.
2. `Pop()`: pop `d`; if `d < 0`, the departing element was the minimum → restore `min = min - d`.
3. `Top()`: peek `d`; if `d < 0` return `min` (the top *is* the minimum), else return `min + d`.
4. `GetMin()`: return `min`.

### Complexity
- **Time:** O(1) per operation — constant arithmetic on the top element and one scalar.
- **Space:** O(n) for the stack itself (unavoidable — the elements must live somewhere), **O(1) auxiliary**: one `int64` per element plus a single scalar, versus two words per element in Approach 2.

### Code
```go
type DiffMinStack struct {
	diffs []int64 // stored as val - minAtPushTime (int64: avoids overflow)
	min   int64   // current minimum of the whole stack
}

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
```

### Dry Run
Example 1:

| Op | diff computed / used | diffs (bottom→top) | min | Returned |
|----|----------------------|--------------------|-----|----------|
| push(-2) | empty → store 0 | [0] | -2 | null |
| push(0) | 0 − (−2) = 2 | [0, 2] | -2 | null |
| push(-3) | −3 − (−2) = **−1** (new low) | [0, 2, −1] | **-3** | null |
| getMin() | — | [0, 2, −1] | -3 | **-3** |
| pop() | d = −1 < 0 → min = −3 − (−1) = **−2** | [0, 2] | **-2** | null |
| top() | d = 2 ≥ 0 → −2 + 2 = 0 | [0, 2] | -2 | **0** |
| getMin() | — | [0, 2] | -2 | **-2** |

Output `[null,null,null,null,-3,null,0,-2]` ✓ — the pop of `−1` resurrected the previous minimum with pure arithmetic.

---

## Key Takeaways

- **Stacks are undo machines**: any "aggregate of the current stack" query (min, max, gcd, …) can be answered in O(1) by snapshotting the aggregate per element — pop automatically rewinds it. This generalizes far beyond min.
- **Lazy auxiliary stack**: store only the values that *change* the aggregate; on pop, compare against the aux top to know when to rewind. Remember `<=` vs `<` for duplicates.
- **Difference encoding** turns "extra storage per element" into "sign information woven into the one number you already store" — a recurring trick for O(1)-auxiliary designs; always check for **overflow** when subtracting two full-range 32-bit values.
- Design-problem interviews reward stating the trade-off table (Approach 2 simplest, 3 leaner in practice, 4 optimal) before writing code.
- The min-stack construction is a building block: two min-stacks make a min-*queue* (see #239 sliding window maximum), and #716 Max Stack extends the same snapshot idea.

---

## Related Problems

- LeetCode #716 — Max Stack (same snapshot idea plus deletion of arbitrary max)
- LeetCode #232 — Implement Queue using Stacks (stack-composition design)
- LeetCode #225 — Implement Stack using Queues (inverse design exercise)
- LeetCode #239 — Sliding Window Maximum (min/max-queue built from two min/max-stacks)
- LeetCode #895 — Maximum Frequency Stack (stack design with per-element bookkeeping)
