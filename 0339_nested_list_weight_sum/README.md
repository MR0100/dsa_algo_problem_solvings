# 0339 — Nested List Weight Sum

> LeetCode #339 · Difficulty: Medium
> **Categories:** Depth-First Search, Breadth-First Search, Stack, Tree

---

## Problem Statement

You are given a nested list of integers `nestedList`. Each element is either an integer or a list whose elements may also be integers or other lists.

The **depth** of an integer is the number of lists that it is inside of. For example, the nested list `[1,[2,2],[[3],2],1]` has each integer's value set to its **depth**. Let `maxDepth` be the **maximum depth** of any integer.

The **weight** of an integer is `depth`.

Return *the sum of each integer in `nestedList` multiplied by its weight (depth)*.

**Example 1:**

```
Input: nestedList = [[1,1],2,[1,1]]
Output: 10
Explanation: Four 1's at depth 2, one 2 at depth 1. 1*2 + 1*2 + 2*1 + 1*2 + 1*2 = 10.
```

**Example 2:**

```
Input: nestedList = [1,[4,[6]]]
Output: 27
Explanation: One 1 at depth 1, one 4 at depth 2, one 6 at depth 3. 1*1 + 4*2 + 6*3 = 27.
```

**Constraints:**

- `1 <= nestedList.length <= 50`
- The values of the integers in the nested list is in the range `[-100, 100]`.
- The maximum depth of any integer is less than or equal to `50`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Facebook   | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| LinkedIn   | ★★★☆☆ Medium     | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Graph BFS/DFS** — the nested list is a tree; both a recursive DFS and a queue-based BFS traverse every node once → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Tree Traversal** — each list is an internal node and each integer a leaf; depth = number of enclosing lists → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Stack** — the iterative variant pushes `(element, depth)` frames to simulate recursion without the call stack → see [`/dsa/stack.md`](/dsa/stack.md)
- **Queue / Deque** — the BFS variant processes the structure one depth-level at a time → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DFS (recursive) | O(N) | O(D) | Cleanest; carry depth as a recursion parameter |
| 2 | BFS (level by level) | O(N) | O(W) | Natural when weight = level; no recursion depth risk |
| 3 | DFS Explicit Stack (Optimal) | O(N) | O(N) | Iterative DFS; avoids call-stack limits |

*N = total elements (ints + lists), D = max depth, W = max width of a level.*

---

## Approach 1 — DFS (recursive)

### Intuition

The weight of an integer is exactly its depth, where the top-level list is depth 1. Recurse through the structure carrying the current depth. A plain integer at depth `d` contributes `value * d`; a nested list is entered at depth `d+1`. Summing all leaf contributions gives the answer.

### Algorithm

1. Define `helper(list, depth)`.
2. For each element: if it's an integer, add `value * depth`; else recurse into its list with `depth+1`.
3. Return `helper(nestedList, 1)`.

### Complexity

- **Time:** O(N) — every element is visited exactly once.
- **Space:** O(D) — recursion stack bounded by the maximum nesting depth.

### Code

```go
func dfs(nestedList []*NestedInteger) int {
	var helper func(list []*NestedInteger, depth int) int
	helper = func(list []*NestedInteger, depth int) int {
		sum := 0
		for _, ni := range list {
			if ni.IsInteger() {
				sum += ni.GetInteger() * depth // leaf contributes value×depth
			} else {
				sum += helper(ni.GetList(), depth+1) // descend: deeper by one level
			}
		}
		return sum
	}
	return helper(nestedList, 1) // top-level list is depth 1
}
```

### Dry Run

Example 1: `nestedList = [[1,1], 2, [1,1]]`, starting `helper(..., 1)`.

| Element | IsInteger? | depth | contribution | running sum |
|---------|-----------|-------|--------------|-------------|
| `[1,1]` | no → recurse depth 2 | — | — | — |
| &nbsp;&nbsp;1 | yes | 2 | 1×2 = 2 | 2 |
| &nbsp;&nbsp;1 | yes | 2 | 1×2 = 2 | 4 |
| `2` | yes | 1 | 2×1 = 2 | 6 |
| `[1,1]` | no → recurse depth 2 | — | — | — |
| &nbsp;&nbsp;1 | yes | 2 | 1×2 = 2 | 8 |
| &nbsp;&nbsp;1 | yes | 2 | 1×2 = 2 | 10 |

Result: `10` ✔

---

## Approach 2 — BFS (level by level)

### Intuition

Since weight equals level, process the structure breadth-first. Hold a queue of the current level's elements. At depth `d`, add each integer's `value * d` and enqueue the children of any list for depth `d+1`. Increment the depth after finishing each level.

### Algorithm

1. `queue = nestedList`, `depth = 1`.
2. While the queue is non-empty: for every element, integers add `value * depth`, lists append their children to the next-level list.
3. Replace the queue with the next-level list, `depth++`.

### Complexity

- **Time:** O(N) — each element enqueued and dequeued once.
- **Space:** O(W) — the queue holds at most one level's worth of elements.

### Code

```go
func bfs(nestedList []*NestedInteger) int {
	queue := nestedList // items pending at the current depth
	depth := 1          // top level is depth 1
	sum := 0
	for len(queue) > 0 {
		var next []*NestedInteger // items collected for the following depth
		for _, ni := range queue {
			if ni.IsInteger() {
				sum += ni.GetInteger() * depth // weigh this integer by its depth
			} else {
				next = append(next, ni.GetList()...) // its children live one level deeper
			}
		}
		queue = next // advance to the next level
		depth++      // and increase the weight
	}
	return sum
}
```

### Dry Run

Example 1: `[[1,1], 2, [1,1]]`.

| depth | queue | integers added | lists → next queue | sum |
|-------|-------|----------------|--------------------|-----|
| 1 | `[1,1]`, `2`, `[1,1]` | 2×1 = 2 | children `1,1,1,1` | 2 |
| 2 | `1,1,1,1` | 1×2 ×4 = 8 | (none) | 10 |
| 3 | (empty) | — | — | 10 |

Result: `10` ✔

---

## Approach 3 — DFS Explicit Stack (Optimal)

### Intuition

Replace recursion with a manual stack of `(element, depth)` frames. Pop a frame: an integer adds `value * depth`; a list pushes each child at `depth + 1`. This is DFS without the call stack — handy when nesting could be deep enough to risk a recursion limit.

### Algorithm

1. Push every top-level element with depth 1.
2. Pop a frame; if integer, add `value * depth`; if list, push each child at `depth + 1`.
3. Repeat until the stack is empty.

### Complexity

- **Time:** O(N) — each element pushed and popped once.
- **Space:** O(N) — worst case the stack holds all elements of a wide level.

### Code

```go
func iterativeStack(nestedList []*NestedInteger) int {
	// frame pairs an element with the depth at which it sits.
	type frame struct {
		ni    *NestedInteger
		depth int
	}
	stack := make([]frame, 0, len(nestedList))
	for _, ni := range nestedList {
		stack = append(stack, frame{ni, 1}) // seed with top-level, depth 1
	}
	sum := 0
	for len(stack) > 0 {
		f := stack[len(stack)-1] // pop the top frame
		stack = stack[:len(stack)-1]
		if f.ni.IsInteger() {
			sum += f.ni.GetInteger() * f.depth // weigh by its depth
		} else {
			for _, child := range f.ni.GetList() {
				stack = append(stack, frame{child, f.depth + 1}) // children go deeper
			}
		}
	}
	return sum
}
```

### Dry Run

Example 1: `[[1,1], 2, [1,1]]`. Seed stack (bottom→top): `{[1,1],1} {2,1} {[1,1],1}`.

| Pop | IsInteger? | action | sum | stack after |
|-----|-----------|--------|-----|-------------|
| `{[1,1],1}` | no | push `{1,2} {1,2}` | 0 | `{[1,1],1} {2,1} {1,2} {1,2}` |
| `{1,2}` | yes | +1×2 | 2 | `{[1,1],1} {2,1} {1,2}` |
| `{1,2}` | yes | +1×2 | 4 | `{[1,1],1} {2,1}` |
| `{2,1}` | yes | +2×1 | 6 | `{[1,1],1}` |
| `{[1,1],1}` | no | push `{1,2} {1,2}` | 6 | `{1,2} {1,2}` |
| `{1,2}` | yes | +1×2 | 8 | `{1,2}` |
| `{1,2}` | yes | +1×2 | 10 | (empty) |

Result: `10` ✔

---

## Key Takeaways

- **A nested list is a tree**: lists are internal nodes, integers are leaves, and "depth" is just the tree level — so any traversal (DFS/BFS/explicit stack) solves it in O(N).
- **Carry the depth with the element**, either as a recursion parameter, a BFS level counter, or a field in a stack frame. All three encode the same weight.
- BFS shines here because *weight equals level*; you never even need to store depth per element — just increment after each level.
- The follow-up **#364 (Nested List Weight Sum II)** inverts the weight (`maxDepth - depth + 1`); a two-pass or "sum-carrying BFS" adapts this cleanly.

---

## Related Problems

- LeetCode #364 — Nested List Weight Sum II (inverted depth weighting)
- LeetCode #341 — Flatten Nested List Iterator (same NestedInteger API, iterator design)
- LeetCode #565 — Array Nesting (different, but "follow the nesting" theme)
- LeetCode #690 — Employee Importance (weighted DFS/BFS over a hierarchy)
