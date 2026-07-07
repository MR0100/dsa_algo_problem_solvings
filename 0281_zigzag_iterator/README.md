# 0281 — Zigzag Iterator

> LeetCode #281 · Difficulty: Medium
> **Categories:** Design, Queue, Iterator, Two Pointers

---

## Problem Statement

Given two vectors of integers `v1` and `v2`, implement an iterator to return their elements alternately.

Implement the `ZigzagIterator` class:

- `ZigzagIterator(List<int> v1, List<int> v2)` initializes the object with the two vectors `v1` and `v2`.
- `boolean hasNext()` returns `true` if the iterator still has elements, and `false` otherwise.
- `int next()` returns the current element of the iterator and moves the iterator to the next element.

**Example 1:**
```
Input: v1 = [1,2], v2 = [3,4,5,6]
Output: [1,3,2,4,5,6]
Explanation: By calling next repeatedly until hasNext returns false, the
order of elements returned by next should be: [1,3,2,4,5,6].
```

**Example 2:**
```
Input: v1 = [1], v2 = []
Output: [1]
```

**Example 3:**
```
Input: v1 = [], v2 = [1]
Output: [1]
```

**Constraints:**
- `0 <= v1.length, v2.length <= 1000`
- `1 <= v1.length + v2.length <= 2000`
- `-2³¹ <= v1[i], v2[i] <= 2³¹ - 1`

**Follow-up:** What if you are given `k` vectors? How well can your code be extended to such cases?

**Clarification for the follow-up:** The "Zigzag" order is not clearly defined and is ambiguous for `k > 2` cases. If "Zigzag" does not look right to you, replace "Zigzag" with "Cyclic". For example:
```
Input: v1 = [1,2,3], v2 = [4,5,6,7], v3 = [8,9]
Output: [1,4,8,2,5,9,3,6,7]
```

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Iterator design** — expose lazy `hasNext`/`next` over a merged view of inputs → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Queue / round-robin scheduling** — a FIFO of live cursors gives the cyclic order and generalises to `k` lists → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Two pointers** — one index per list with a turn toggle merges lazily in O(1) space → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview
| # | Approach | Time (next) | Space | When to use |
|---|----------|-------------|-------|-------------|
| 1 | Flatten Upfront | O(1) | O(n+m) | Simplest; fine when full materialisation is acceptable |
| 2 | Two Cursors + Turn Toggle | O(1) | O(1) | Lazy, minimal memory, exactly two lists |
| 3 | Queue of Iterators (Optimal) | O(1) | O(k) | Generalises to k lists cleanly |

---

## Approach 1 — Flatten Upfront

### Intuition
The output is just the two lists interleaved column by column: `v1[0], v2[0], v1[1], v2[1], …`, and when one list is exhausted you keep draining the other. Build that merged order once at construction, then iterating is a plain slice walk.

### Algorithm
1. Loop `j = 0, 1, 2, …` until `j` is past both lists.
2. At each `j`, append `v1[j]` if it exists, then `v2[j]` if it exists.
3. `hasNext` = cursor `< len(merged)`; `next` returns `merged[pos]` and advances `pos`.

### Complexity
- **Time:** Constructor O(n+m); `next`/`hasNext` O(1) — one slice read.
- **Space:** O(n+m) — the flattened buffer holds every element.

### Code
```go
type FlattenZigzag struct {
	merged []int // both lists interleaved in zigzag order
	pos    int   // index of the next element to emit
}

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

func (z *FlattenZigzag) HasNext() bool { return z.pos < len(z.merged) }

func (z *FlattenZigzag) Next() int {
	v := z.merged[z.pos]
	z.pos++
	return v
}
```

### Dry Run
`v1 = [1,2]`, `v2 = [3,4,5,6]`.

| j | append v1[j] | append v2[j] | merged so far |
|---|--------------|--------------|---------------|
| 0 | 1            | 3            | [1,3]         |
| 1 | 2            | 4            | [1,3,2,4]     |
| 2 | —            | 5            | [1,3,2,4,5]   |
| 3 | —            | 6            | [1,3,2,4,5,6] |

Iterating `merged` yields `[1,3,2,4,5,6]`.

---

## Approach 2 — Two Cursors, Turn Toggle

### Intuition
No need to buffer everything. Keep a pointer into each list and a `turn` flag saying whose element comes next. Before emitting, if the preferred list is exhausted, fall back to the other, so the chosen list always has data. This is O(1) extra space and fully lazy.

### Algorithm
1. State: `i1`, `i2` (per-list cursors), `turn ∈ {0,1}`.
2. `hasNext`: `i1 < len(v1) || i2 < len(v2)`.
3. `next`: if it's `v1`'s turn and `v1` still has data, or `v2` is exhausted → take `v1[i1++]`, set `turn = 1`; else take `v2[i2++]`, set `turn = 0`.

### Complexity
- **Time:** `next` O(1) amortised; `hasNext` O(1).
- **Space:** O(1) — only two indices and a flag.

### Code
```go
type TwoCursorZigzag struct {
	v1, v2 []int
	i1, i2 int
	turn   int
}

func NewTwoCursorZigzag(v1, v2 []int) *TwoCursorZigzag {
	return &TwoCursorZigzag{v1: v1, v2: v2}
}

func (z *TwoCursorZigzag) HasNext() bool {
	return z.i1 < len(z.v1) || z.i2 < len(z.v2)
}

func (z *TwoCursorZigzag) Next() int {
	if (z.turn == 0 && z.i1 < len(z.v1)) || z.i2 >= len(z.v2) {
		v := z.v1[z.i1]
		z.i1++
		z.turn = 1
		return v
	}
	v := z.v2[z.i2]
	z.i2++
	z.turn = 0
	return v
}
```

### Dry Run
`v1 = [1,2]`, `v2 = [3,4,5,6]`.

| Call | turn (in) | branch chosen | emitted | i1,i2 (out) | turn (out) |
|------|-----------|---------------|---------|-------------|------------|
| 1 | 0 | v1 (turn 0, i1<2) | 1 | 1,0 | 1 |
| 2 | 1 | v2 | 3 | 1,1 | 0 |
| 3 | 0 | v1 | 2 | 2,1 | 1 |
| 4 | 1 | v2 | 4 | 2,2 | 0 |
| 5 | 0 | v1 empty → v2 (i2<6) | 5 | 2,3 | 0 |
| 6 | 0 | v1 empty → v2 | 6 | 2,4 | 0 |

Output: `[1,3,2,4,5,6]`.

---

## Approach 3 — Queue of Iterators (Optimal)

### Intuition
Model each list as a cursor `(slice, index)`. Keep a FIFO queue of the cursors that still have elements. To emit, pop the front cursor, take its current value, advance it, and push it back **only if** it still has data. The queue naturally cycles `list1 → list2 → … → listk → list1`, i.e. the zigzag/cyclic order, and drops lists as they empty — so it extends to `k` lists for free.

### Algorithm
1. Constructor: enqueue a cursor for every non-empty list.
2. `hasNext`: queue is non-empty.
3. `next`: dequeue cursor `c`; `value = c.list[c.idx]`; `c.idx++`; if `c` still has data, enqueue it again; return `value`.

### Complexity
- **Time:** `next` O(1); `hasNext` O(1).
- **Space:** O(k) — at most one cursor per list in the queue.

### Code
```go
type cursor struct {
	list []int
	idx  int
}

type QueueZigzag struct {
	queue []*cursor
}

func NewQueueZigzag(lists ...[]int) *QueueZigzag {
	z := &QueueZigzag{}
	for _, l := range lists {
		if len(l) > 0 {
			z.queue = append(z.queue, &cursor{list: l})
		}
	}
	return z
}

func (z *QueueZigzag) HasNext() bool { return len(z.queue) > 0 }

func (z *QueueZigzag) Next() int {
	c := z.queue[0]
	z.queue = z.queue[1:]
	v := c.list[c.idx]
	c.idx++
	if c.idx < len(c.list) {
		z.queue = append(z.queue, c)
	}
	return v
}
```

### Dry Run
`v1 = [1,2]` (cursor A), `v2 = [3,4,5,6]` (cursor B). Queue starts `[A(0), B(0)]`.

| Call | front | value | advance | requeue? | queue after |
|------|-------|-------|---------|----------|-------------|
| 1 | A(0) | 1 | A→1 | yes | [B(0), A(1)] |
| 2 | B(0) | 3 | B→1 | yes | [A(1), B(1)] |
| 3 | A(1) | 2 | A→2 (done) | no | [B(1)] |
| 4 | B(1) | 4 | B→2 | yes | [B(2)] |
| 5 | B(2) | 5 | B→3 | yes | [B(3)] |
| 6 | B(3) | 6 | B→4 (done) | no | [] |

Output: `[1,3,2,4,5,6]`.

---

## Key Takeaways
- **Round-robin via a queue** is the cleanest way to interleave `k` sequences and answers the follow-up with zero code changes.
- Interleaving iterators is a **merge** problem: choose eager materialisation (simple) vs. lazy cursors (O(1) space) based on whether early termination and memory matter.
- The "skip exhausted list" guard is what keeps the turn logic correct once one list runs out.

---

## Related Problems
- LeetCode #341 — Flatten Nested List Iterator (iterator design over nested data)
- LeetCode #284 — Peeking Iterator (wrapping an iterator with extra state)
- LeetCode #251 — Flatten 2D Vector (row cursors over a 2D structure)
- LeetCode #23 — Merge k Sorted Lists (k-way merge, priority instead of round-robin)
