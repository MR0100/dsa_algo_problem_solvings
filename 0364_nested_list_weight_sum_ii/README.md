# 0364 — Nested List Weight Sum II

> LeetCode #364 · Difficulty: Medium
> **Categories:** Depth-First Search, Breadth-First Search, Stack, Tree

---

## Problem Statement

You are given a nested list of integers `nestedList`. Each element is either an integer or a list whose elements may also be integers or other lists.

The **depth** of an integer is the number of lists that it is inside of. For example, the nested list `[1,[2,2],[[3],2],1]` has each integer's value set to its depth. Let `maxDepth` be the **maximum depth** of any integer.

The **weight** of an integer is `maxDepth - (the depth of the integer) + 1`.

Return *the sum of each integer in* `nestedList` *multiplied by its weight*.

**Example 1:**

```
Input: nestedList = [[1,1],2,[1,1]]
Output: 8
Explanation: Four 1's with a weight of 1, one 2 with a weight of 2.
1*1 + 1*1 + 2*2 + 1*1 + 1*1 = 8
```

**Example 2:**

```
Input: nestedList = [1,[4,[6]]]
Output: 17
Explanation: One 1 at depth 3, one 4 at depth 2, and one 6 at depth 1.
1*3 + 4*2 + 6*1 = 17
```

**Constraints:**

- `1 <= nestedList.length <= 50`
- The values of the integers in the nested list is in the range `[-100, 100]`.
- The maximum depth of any integer is less than or equal to `50`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| LinkedIn   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Depth-First Search** — recursion over the nested structure computes depths and weighted sums naturally → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Breadth-First Search** — the one-pass optimal walks level by level, accumulating running level sums → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Tree Traversal** — a nested list is a tree; integers are leaves, lists are internal nodes → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two-Pass DFS | O(N) | O(D) | Most direct; find maxDepth, then weight bottom-up |
| 2 | One-Pass BFS (Optimal) | O(N) | O(W) | Single traversal, no depth pre-pass; elegant accumulation |

*(N = total nodes, D = max depth, W = max width of one level.)*

---

## Approach 1 — Two-Pass DFS

### Intuition

Weight here is **inverted** vs. problem 339: leaves weigh 1 and shallow integers weigh most. If `maxDepth` is the deepest nesting level, an integer at depth `d` has weight `maxDepth - d + 1`. So first find `maxDepth` with one DFS, then a second DFS multiplies each integer by its inverted weight.

### Algorithm

1. DFS to compute `maxDepth` (top-level items are depth 1).
2. DFS again; for each integer at depth `d`, add `value * (maxDepth - d + 1)`.
3. Return the total.

### Complexity

- **Time:** O(N) — two linear passes over all nodes.
- **Space:** O(D) — recursion stack bounded by nesting depth.

### Code

```go
func twoPassDFS(nestedList []*NestedInteger) int {
	maxDepth := findMaxDepth(nestedList, 1) // depth of top-level items is 1
	return weightedSum(nestedList, 1, maxDepth)
}

func findMaxDepth(list []*NestedInteger, depth int) int {
	best := depth // at least this deep just by being here
	for _, ni := range list {
		if !ni.IsInteger() {
			// Recurse one level deeper into the sublist.
			if d := findMaxDepth(ni.GetList(), depth+1); d > best {
				best = d
			}
		}
	}
	return best
}

func weightedSum(list []*NestedInteger, depth, maxDepth int) int {
	sum := 0
	for _, ni := range list {
		if ni.IsInteger() {
			// Weight is inverted: shallow items weigh more, leaves weigh 1.
			sum += ni.GetInteger() * (maxDepth - depth + 1)
		} else {
			sum += weightedSum(ni.GetList(), depth+1, maxDepth)
		}
	}
	return sum
}
```

### Dry Run

Example 1: `[[1,1],2,[1,1]]`.

Pass 1 — `findMaxDepth`:

| Element | depth | contributes |
|---------|-------|-------------|
| `[1,1]` | recurse → ints at depth 2 | 2 |
| `2`     | depth 1 | 1 |
| `[1,1]` | recurse → depth 2 | 2 |

`maxDepth = 2`.

Pass 2 — `weightedSum` with `maxDepth=2`:

| Integer | depth d | weight = 2-d+1 | contribution |
|---------|---------|----------------|--------------|
| 1 (in first list) | 2 | 1 | 1 |
| 1 (in first list) | 2 | 1 | 1 |
| 2                 | 1 | 2 | 4 |
| 1 (in third list) | 2 | 1 | 1 |
| 1 (in third list) | 2 | 1 | 1 |

Total = 1+1+4+1+1 = `8` ✔

---

## Approach 2 — One-Pass BFS (Optimal)

### Intuition

Keep a running `levelSum` of every integer seen down to the current level, and add `levelSum` into `total` once per level. An integer that first appears at level `ℓ` is thereby counted for level `ℓ`, `ℓ+1`, …, down to the last level — exactly `maxDepth - ℓ + 1` times, which is its weight. So a single BFS accumulates the inverted-weight sum with no separate depth pass.

### Algorithm

1. Queue = the top-level items. `levelSum = 0`, `total = 0`.
2. For each level: add every integer at this level into `levelSum`; collect children of every list for the next level. Then `total += levelSum`.
3. When the queue empties, return `total`.

### Complexity

- **Time:** O(N) — each node enqueued and dequeued once.
- **Space:** O(W) — the frontier holds at most one level's worth of nodes.

### Code

```go
func onePassBFS(nestedList []*NestedInteger) int {
	queue := append([]*NestedInteger{}, nestedList...) // level 1 frontier
	levelSum, total := 0, 0
	for len(queue) > 0 {
		var next []*NestedInteger // nodes forming the next level
		for _, ni := range queue {
			if ni.IsInteger() {
				levelSum += ni.GetInteger() // carried into every deeper level
			} else {
				next = append(next, ni.GetList()...) // descend one level
			}
		}
		// Adding levelSum once per remaining level accumulates each integer's
		// inverted weight (maxDepth - itsDepth + 1) automatically.
		total += levelSum
		queue = next
	}
	return total
}
```

### Dry Run

Example 1: `[[1,1],2,[1,1]]`.

| Level | queue contents | integers added to levelSum | levelSum after | total += levelSum | total |
|-------|----------------|----------------------------|----------------|-------------------|-------|
| 1 | `[1,1]`, `2`, `[1,1]` | just the `2` | 2 | +2 | 2 |
| 2 | `1,1,1,1` (children of both lists) | 1+1+1+1 = 4 | 2+4 = 6 | +6 | 8 |
| — | empty | — | — | stop | **8** |

Total = `8` ✔ — the `2` (depth 1) is added at both levels → counted twice (weight 2); the four `1`s (depth 2) are added only at level 2 → counted once (weight 1).

---

## Key Takeaways

- **"Weight grows toward the root" = inverted weight** `maxDepth - depth + 1`. This is the twist that separates #364 from #339 (where weight grows with depth).
- **Two clean strategies:** (a) find `maxDepth` first, then weight; or (b) a running-level-sum BFS where re-adding the accumulated sum each level *is* the inverted weight — no depth needed up front.
- The BFS trick — "add the prefix of level sums" — is a neat way to turn "multiply by (maxDepth − depth + 1)" into repeated addition, avoiding a second pass.
- A nested list is just a tree: integers are leaves, lists are internal nodes; both DFS and BFS templates apply directly.

---

## Related Problems

- LeetCode #339 — Nested List Weight Sum (weight = depth, not inverted)
- LeetCode #341 — Flatten Nested List Iterator (same NestedInteger interface)
- LeetCode #565 — Array Nesting (nested structure traversal)
- LeetCode #690 — Employee Importance (weighted BFS/DFS over a tree)
