# 0270 — Closest Binary Search Tree Value

> LeetCode #270 · Difficulty: Easy
> **Categories:** Binary Search, Tree, Depth-First Search, Binary Search Tree

---

## Problem Statement

Given the `root` of a binary search tree and a `target` value, return the value
in the BST that is closest to the `target`. If there are multiple answers, print
the **smallest**.

**Example 1:**
```
Input: root = [4,2,5,1,3], target = 3.714286
Output: 4
```

**Example 2:**
```
Input: root = [1], target = 4.428571
Output: 1
```

**Constraints:**
- The number of nodes in the tree is in the range `[1, 10^4]`.
- `0 <= Node.val <= 10^9`
- `-10^9 <= target <= 10^9`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Binary Search Tree** — the BST ordering steers the search toward the target → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **Binary Search** — each step discards one subtree, like halving a search space → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Tree Traversal** — in-order walk flattens the tree to a sorted sequence → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | In-order Traversal + Scan | O(n) | O(n) | Simple, ignores BST shape |
| 2 | Iterative BST Descent (Optimal) | O(h) | O(1) | Exploit BST ordering, no extra memory |

> `h` = tree height: O(log n) balanced, O(n) skewed.

---

## Approach 1 — In-order Traversal + Linear Scan

### Intuition
Ignore the tree shape entirely: collect every node's value, then pick the one
with the smallest absolute distance to `target`. Always correct, easy to reason
about. Scanning in ascending (in-order) sequence and replacing only on a
**strictly** smaller distance naturally satisfies the "smallest on ties" rule.

### Algorithm
1. In-order traverse to collect all values.
2. Track the value with minimum `|val - target|`.
3. On ties, prefer the smaller value (guaranteed by scanning ascending and
   replacing only on `<`).

### Complexity
- **Time:** O(n) — visits every node.
- **Space:** O(n) — the collected slice plus recursion stack.

### Code
```go
func inorderScan(root *TreeNode, target float64) int {
	var vals []int
	var walk func(n *TreeNode)
	walk = func(n *TreeNode) {
		if n == nil {
			return
		}
		walk(n.Left)               // left subtree (smaller values)
		vals = append(vals, n.Val) // visit node
		walk(n.Right)              // right subtree (larger values)
	}
	walk(root)

	closest := vals[0]
	best := math.Abs(float64(closest) - target) // smallest distance so far
	for _, v := range vals[1:] {
		d := math.Abs(float64(v) - target)
		// strictly-smaller distance wins; equal distance keeps the smaller
		// value because we scan in ascending order and only replace on <.
		if d < best {
			best = d
			closest = v
		}
	}
	return closest
}
```

### Dry Run
Tree `[4,2,5,1,3]`, `target = 3.714286`. In-order values: `[1,2,3,4,5]`.

| v | \|v − 3.714286\| | best so far | closest |
|---|------------------|-------------|---------|
| 1 | 2.714            | 2.714       | 1       |
| 2 | 1.714            | 1.714       | 2       |
| 3 | 0.714            | 0.714       | 3       |
| 4 | 0.286            | 0.286       | 4       |
| 5 | 1.286            | 0.286       | 4       |

Return **4**. ✅

---

## Approach 2 — Iterative BST Descent (Optimal)

### Intuition
At each node the BST property tells us which half of the values lie ahead: go
**left** if `target < node.Val`, otherwise **right**. The single root-to-leaf
path we follow passes every candidate that could be closest, so tracking the
best along the way suffices — we never need to see the whole tree.

### Algorithm
1. `closest = root.Val`.
2. While `node != nil`: update `closest` if this node is nearer (ties → smaller
   value). Then move left if `target < node.Val`, else right.
3. Return `closest`.

### Complexity
- **Time:** O(h) — one node per level.
- **Space:** O(1) — no recursion, no extra storage.

### Code
```go
func bstDescent(root *TreeNode, target float64) int {
	closest := root.Val
	node := root
	for node != nil {
		// Prefer this node if strictly closer, or equally close but smaller
		// (LeetCode breaks ties toward the smaller value).
		curD := math.Abs(float64(node.Val) - target)
		bestD := math.Abs(float64(closest) - target)
		if curD < bestD || (curD == bestD && node.Val < closest) {
			closest = node.Val
		}
		if target < float64(node.Val) {
			node = node.Left // target is smaller -> smaller values are left
		} else {
			node = node.Right // target is larger -> larger values are right
		}
	}
	return closest
}
```

### Dry Run
Tree `[4,2,5,1,3]`, `target = 3.714286`:

| node | val | \|val−t\| | closest after | move (target<val?) |
|------|-----|-----------|---------------|--------------------|
| root | 4   | 0.286     | 4             | 3.714 < 4 → left   |
| 2    | 2   | 1.714     | 4 (0.286 wins)| 3.714 ≥ 2 → right  |
| 3    | 3   | 0.714     | 4 (0.286 wins)| 3.714 ≥ 3 → right  |
| nil  | —   | —         | 4             | stop               |

Return **4**. ✅

---

## Key Takeaways
- The BST invariant lets you **discard an entire subtree** at each step — a binary search over tree values.
- The optimal walk is O(h) space-free; recursion or a values array is unnecessary.
- Handle the **tie rule** ("smallest value") explicitly: replace on strictly-closer, or equal-distance-but-smaller.
- Float comparison against integer node values needs a cast; `math.Abs` on `float64` gives the distance.

---

## Related Problems
- LeetCode #272 — Closest Binary Search Tree Value II (k closest values)
- LeetCode #700 — Search in a Binary Search Tree (same descent skeleton)
- LeetCode #235 — Lowest Common Ancestor of a BST (BST-directed descent)
- LeetCode #98 — Validate Binary Search Tree (BST invariant)
