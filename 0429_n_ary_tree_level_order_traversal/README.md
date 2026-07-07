# 0429 — N-ary Tree Level Order Traversal

> LeetCode #429 · Difficulty: Medium
> **Categories:** Tree, Breadth-First Search, Depth-First Search

---

## Problem Statement

Given an n-ary tree, return the *level order* traversal of its nodes' values.

*Nary-Tree input serialization is represented in their level order traversal, each group of children is separated by the null value (See examples).*

**Example 1:**
```
Input: root = [1,null,3,2,4,null,5,6]
Output: [[1],[3,2,4],[5,6]]
```
The tree looks like:
```
        1
      / | \
     3  2  4
    / \
   5   6
```

**Example 2:**
```
Input: root = [1,null,2,3,4,5,null,null,6,7,null,8,null,9,10,null,null,11,null,12,null,13,null,null,14]
Output: [[1],[2,3,4,5],[6,7,8,9,10],[11,12,13],[14]]
```

**Constraints:**
- The height of the n-ary tree is less than or equal to `1000`.
- The total number of nodes is between `[0, 10⁴]`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Breadth-First Search (level-by-level)** — a FIFO queue drained one level at a time is the canonical way to bucket nodes by depth → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Tree Traversal** — the N-ary generalisation of binary-tree level order; children are visited left-to-right → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Queue / Deque** — the ordered container whose FIFO property is exactly what "left to right, level by level" requires → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

Let n = number of nodes, h = tree height.

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS with an explicit queue (Optimal) | O(n) | O(n) | The natural, interview-standard answer; iterative, no recursion depth risk |
| 2 | DFS carrying the depth | O(n) | O(n) (O(h) stack) | When you already have a recursive skeleton, or want the values grouped by a depth index |

---

## Approach 1 — BFS with an Explicit Queue (Optimal)

### Intuition
Level order *is* breadth-first search. A FIFO queue naturally visits nodes in left-to-right, top-to-bottom order — but a plain BFS loses the level boundaries we must report. The fix is a one-line trick: **before draining, snapshot the queue length.** At the instant a level starts, the queue holds *exactly* that level's nodes (their children haven't been enqueued yet), so that count tells us how many pops belong to the current bucket. Pop exactly that many, collect their values, and push their children to seed the next level.

### Algorithm
1. If `root` is `nil`, return `[]`.
2. Seed a queue with `root`.
3. While the queue is non-empty:
   1. `width = len(queue)` — the number of nodes on the current level.
   2. Pop `width` nodes; append each `Val` to a `level` slice; enqueue each node's `Children` (already left-to-right).
   3. Append `level` to the answer.

### Complexity
- **Time:** O(n) — each node is enqueued once and dequeued once; appending its children is O(children).
- **Space:** O(n) — the queue holds at most one level (which can be O(n) wide), plus the O(n) output.

### Code
```go
func bfs(root *Node) [][]int {
	result := [][]int{} // final list of levels; stays [] for an empty tree
	if root == nil {
		return result // nothing to traverse
	}

	queue := []*Node{root} // FIFO seeded with the single root node
	for len(queue) > 0 {
		width := len(queue)   // # nodes on this level = everything currently queued
		level := []int{}      // values collected for this level, left to right
		for i := 0; i < width; i++ {
			node := queue[0]  // front of the queue
			queue = queue[1:] // dequeue (advance the head)
			level = append(level, node.Val)
			// Children are already stored left-to-right, so enqueueing them in
			// order preserves the left-to-right requirement on the next level.
			queue = append(queue, node.Children...)
		}
		result = append(result, level) // this whole level is done
	}
	return result
}
```

### Dry Run (Example 1)

Tree: `1 → (3,2,4)`, `3 → (5,6)`.

| Step | `width` | Nodes popped | `level` built | Children enqueued | `queue` after | `result` |
|------|---------|--------------|---------------|-------------------|---------------|----------|
| seed | — | — | — | — | `[1]` | `[]` |
| 1 | 1 | 1 | `[1]` | 3,2,4 | `[3,2,4]` | `[[1]]` |
| 2 | 3 | 3,2,4 | `[3,2,4]` | 5,6 (from 3); none from 2,4 | `[5,6]` | `[[1],[3,2,4]]` |
| 3 | 2 | 5,6 | `[5,6]` | none | `[]` | `[[1],[3,2,4],[5,6]]` |

Queue empty → output `[[1],[3,2,4],[5,6]]` ✓

---

## Approach 2 — DFS Carrying the Depth

### Intuition
A node at depth `d` belongs in `result[d]`, *no matter when it is visited*. DFS and BFS disagree on visit **order** but never on which **level** a node occupies. So walk the tree depth-first, threading the current depth through the recursion, and append each value into the bucket for its depth. The only bookkeeping: the first time we descend to a new depth, `result` has no slice there yet, so we open a fresh one.

### Algorithm
1. Recurse from `root` at depth `0`.
2. On entering a node at depth `d`: if `d == len(result)`, append a new empty slice (first node seen at this depth).
3. Append the node's value to `result[d]`.
4. Recurse into each child at depth `d+1`, left to right.

### Complexity
- **Time:** O(n) — every node is visited exactly once.
- **Space:** O(n) — recursion stack is O(h) (up to O(n) for a degenerate chain), plus the O(n) output.

### Code
```go
func dfs(root *Node) [][]int {
	result := [][]int{}
	var walk func(node *Node, depth int)
	walk = func(node *Node, depth int) {
		if node == nil {
			return
		}
		if depth == len(result) {
			// First time we descend to this depth: open a new level bucket.
			result = append(result, []int{})
		}
		result[depth] = append(result[depth], node.Val) // place value on its level
		for _, c := range node.Children {
			walk(c, depth+1) // children live one level deeper
		}
	}
	walk(root, 0)
	return result
}
```

### Dry Run (Example 1)

`walk(node, depth)` calls in order; `+bucket` marks a new level being opened:

| Call | `depth` | `len(result)` before | Action | `result` after |
|------|---------|----------------------|--------|----------------|
| `walk(1,0)` | 0 | 0 | +bucket, append 1 | `[[1]]` |
| `walk(3,1)` | 1 | 1 | +bucket, append 3 | `[[1],[3]]` |
| `walk(5,2)` | 2 | 2 | +bucket, append 5 | `[[1],[3],[5]]` |
| `walk(6,2)` | 2 | 3 | append 6 | `[[1],[3],[5,6]]` |
| `walk(2,1)` | 1 | 3 | append 2 | `[[1],[3,2],[5,6]]` |
| `walk(4,1)` | 1 | 3 | append 4 | `[[1],[3,2,4],[5,6]]` |

Output `[[1],[3,2,4],[5,6]]` ✓ — note DFS fills `[5,6]` *before* finishing level 1, yet the depth index keeps every value in the right bucket.

---

## Key Takeaways

- **"Snapshot the queue length" is the universal level-order idiom.** It works identically for binary trees (#102), N-ary trees (this problem), and grid BFS with layers — the queue at the start of an iteration always holds one complete frontier.
- **BFS and DFS both produce level order** as long as DFS is given the depth: the level a node lands on is a property of the node, independent of traversal order.
- N-ary vs binary is a trivial change: instead of pushing `left` then `right`, push all of `Children` — already stored left-to-right, so ordering is free.
- Initialise the answer to `[][]int{}` (not `nil`) so the empty-tree case returns `[]` rather than `null`.

---

## Related Problems
- LeetCode #102 — Binary Tree Level Order Traversal (the binary-tree original of this pattern)
- LeetCode #107 — Binary Tree Level Order Traversal II (same, bottom-up)
- LeetCode #103 — Binary Tree Zigzag Level Order Traversal (alternate row direction)
- LeetCode #559 — Maximum Depth of N-ary Tree (same N-ary traversal, tracking depth)
- LeetCode #589 — N-ary Tree Preorder Traversal (same tree, DFS order)
- LeetCode #590 — N-ary Tree Postorder Traversal (same tree, DFS order)
