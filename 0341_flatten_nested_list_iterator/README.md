# 0341 — Flatten Nested List Iterator

> LeetCode #341 · Difficulty: Medium
> **Categories:** Stack, Tree, Depth-First Search, Design, Queue, Iterator

---

## Problem Statement

You are given a nested list of integers `nestedList`. Each element is either an integer or a list whose elements may also be integers or other lists. Implement an iterator to flatten it.

Implement the `NestedIterator` class:

- `NestedIterator(List<NestedInteger> nestedList)` Initializes the iterator with the nested list `nestedList`.
- `int next()` Returns the next integer in the nested list.
- `boolean hasNext()` Returns `true` if there are still some integers in the nested list and `false` otherwise.

Your code will be tested with the following pseudocode:

```
initialize iterator with nestedList
res = []
while iterator.hasNext()
    append iterator.next() to the end of res
return res
```

If `res` matches the expected flattened list, then your code will be judged as correct.

**Example 1:**

```
Input: nestedList = [[1,1],2,[1,1]]
Output: [1,1,2,1,1]
Explanation: By calling next repeatedly until hasNext returns false,
the order of elements returned by next should be: [1,1,2,1,1].
```

**Example 2:**

```
Input: nestedList = [1,[4,[6]]]
Output: [1,4,6]
Explanation: By calling next repeatedly until hasNext returns false,
the order of elements returned by next should be: [1,4,6].
```

**Constraints:**

- `1 <= nestedList.length <= 500`
- The values of the integers in the nested list is in the range `[-10^6, 10^6]`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Design / Iterator** — the task is to implement a class exposing `hasNext`/`next` over a recursive structure; the interviewer usually wants the *lazy* version → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Stack** — an explicit stack replaces the recursion call stack, letting us unpack lists on demand instead of all at once → see [`/dsa/stack.md`](/dsa/stack.md)
- **Depth-First Search** — the flatten order is a left-to-right DFS of the nested tree → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Eager Flatten in Constructor | O(N) build, O(1)/op | O(N) | Simple to write; fine when all integers fit in memory |
| 2 | Lazy Stack (Optimal) | O(1) amortised/op | O(D) depth | Interview-preferred; memory scales with nesting depth, not size |

---

## Approach 1 — Eager Flatten in Constructor

### Intuition
The required output order is exactly a depth-first, left-to-right walk of the nested tree. Do that walk once in the constructor, recording every integer into a flat slice; then `next`/`hasNext` are trivial index operations.

### Algorithm
1. In the constructor, DFS the nested list. Append leaf integers; recurse into sub-lists.
2. `HasNext`: `pos < len(flat)`.
3. `Next`: return `flat[pos]`, then `pos++`.

### Complexity
- **Time:** O(N) to build (N = total integers); O(1) per `Next`/`HasNext`.
- **Space:** O(N) — the flat slice stores every integer.

### Code
```go
type NestedIteratorEager struct {
	flat []int
	pos  int
}

func ConstructorEager(nestedList []*NestedInteger) *NestedIteratorEager {
	it := &NestedIteratorEager{}
	var dfs func(list []*NestedInteger)
	dfs = func(list []*NestedInteger) {
		for _, ni := range list {
			if ni.IsInteger() {
				it.flat = append(it.flat, ni.GetInteger())
			} else {
				dfs(ni.GetList())
			}
		}
	}
	dfs(nestedList)
	return it
}

func (it *NestedIteratorEager) HasNext() bool { return it.pos < len(it.flat) }

func (it *NestedIteratorEager) Next() int {
	v := it.flat[it.pos]
	it.pos++
	return v
}
```

### Dry Run
Input `[[1,1],2,[1,1]]`, constructor DFS:

| Step | Element visited | Action | flat |
|------|-----------------|--------|------|
| 1 | `[1,1]` | list → recurse | `[]` |
| 2 | `1` | leaf → append | `[1]` |
| 3 | `1` | leaf → append | `[1,1]` |
| 4 | `2` | leaf → append | `[1,1,2]` |
| 5 | `[1,1]` | list → recurse | `[1,1,2]` |
| 6 | `1` | leaf → append | `[1,1,2,1]` |
| 7 | `1` | leaf → append | `[1,1,2,1,1]` |

Then `Next` × 5 returns `1,1,2,1,1`. Output `[1,1,2,1,1]`.

---

## Approach 2 — Lazy Stack (Optimal)

### Intuition
Recursion uses the call stack; make it explicit. Keep a stack of pending `NestedInteger`s. Only unpack a list when we actually reach it, so at any moment memory holds just the path being explored — proportional to nesting depth, not total size.

### Algorithm
1. Constructor: push the top-level items onto the stack **in reverse** so the first logical element is on top.
2. `HasNext`: while the top is a list, pop it and push its children reversed; return `true` iff the stack is non-empty afterwards (top is now an integer).
3. `Next`: call `HasNext` to prime the stack, then pop and return the top integer.

### Complexity
- **Time:** O(1) amortised per operation — each `NestedInteger` node is pushed and popped exactly once across the whole iteration.
- **Space:** O(D) — the stack never holds more than the current nesting depth plus sibling lists.

### Code
```go
type NestedIteratorStack struct {
	stack []*NestedInteger
}

func ConstructorStack(nestedList []*NestedInteger) *NestedIteratorStack {
	it := &NestedIteratorStack{}
	for i := len(nestedList) - 1; i >= 0; i-- {
		it.stack = append(it.stack, nestedList[i])
	}
	return it
}

func (it *NestedIteratorStack) HasNext() bool {
	for len(it.stack) > 0 {
		top := it.stack[len(it.stack)-1]
		if top.IsInteger() {
			return true
		}
		it.stack = it.stack[:len(it.stack)-1]
		list := top.GetList()
		for i := len(list) - 1; i >= 0; i-- {
			it.stack = append(it.stack, list[i])
		}
	}
	return false
}

func (it *NestedIteratorStack) Next() int {
	it.HasNext()
	top := it.stack[len(it.stack)-1]
	it.stack = it.stack[:len(it.stack)-1]
	return top.GetInteger()
}
```

### Dry Run
Input `[[1,1],2,[1,1]]`. Constructor pushes reversed: stack (top on right) = `[ [1,1], 2, [1,1] ]` where the top is the first `[1,1]`... shown reversed below. Top-of-stack is rightmost.

| Call | Stack before (top→right) | Action | Returns |
|------|--------------------------|--------|---------|
| build | `[1,1] , 2 , [1,1]` (reversed) | push reversed → top = first `[1,1]` | — |
| HasNext | top=`[1,1]` | list → pop, push `1,1` reversed; top=`1` | true |
| Next | top=`1` | pop | 1 |
| HasNext | top=`1` | integer | true |
| Next | top=`1` | pop | 1 |
| HasNext | top=`2` | integer | true |
| Next | top=`2` | pop | 2 |
| HasNext | top=`[1,1]` | list → pop, push `1,1`; top=`1` | true |
| Next | top=`1` | pop | 1 |
| Next | (primed) | pop | 1 |
| HasNext | empty | — | false |

Output `[1,1,2,1,1]`.

---

## Key Takeaways

- The "flatten an iterator" pattern is really "turn recursion into an explicit stack, then make it lazy."
- Push children **reversed** so left-to-right order is preserved when popping from the top.
- Put the list-unpacking logic in `HasNext` and have `Next` call `HasNext` first — this keeps the top of the stack always primed to a real integer and avoids duplicating logic.
- Eager pre-flattening is simpler and perfectly acceptable when the whole structure fits in memory; the lazy stack wins when inputs are huge or partially consumed.

---

## Related Problems

- LeetCode #251 — Flatten 2D Vector (same lazy-iterator pattern, one level)
- LeetCode #281 — Zigzag Iterator (multi-source iterator design)
- LeetCode #385 — Mini Parser (parsing into the same NestedInteger structure)
- LeetCode #173 — Binary Search Tree Iterator (stack-based lazy iterator over a tree)
- LeetCode #339 — Nested List Weight Sum (DFS over the same nested structure)
