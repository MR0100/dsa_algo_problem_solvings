# 0298 — Binary Tree Longest Consecutive Sequence

> LeetCode #298 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return *the length of the longest consecutive sequence path*.

A **consecutive sequence path** is a path where the values **increase by one** along the path.

Note that the path can start **at any node** in the tree, and you **cannot** go from a node to its parent in the path.

**Example 1:**

```
Input: root = [1,null,3,2,4,null,null,null,5]
Output: 3
Explanation: Longest consecutive sequence path is 3-4-5, so return 3.
```

**Example 2:**

```
Input: root = [2,null,3,2,null,1]
Output: 2
Explanation: Longest consecutive sequence path is 2-3, not 3-2-1, so return 2.
```

**Constraints:**

- The number of nodes in the tree is in the range `[1, 3 * 10^4]`.
- `-3 * 10^4 <= Node.val <= 3 * 10^4`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Depth-First Search on trees** — the whole solution is one DFS pass → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Top-down vs bottom-up recursion** — carrying state down vs returning it up, a core recursion pattern → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Top-Down DFS (pass length down) | O(n) | O(h) | Most intuitive; carries the run length as a parameter |
| 2 | Bottom-Up DFS (return length up) | O(n) | O(h) | Clean postorder; no extra parameters |

*(n = nodes, h = height.)*

---

## Approach 1 — Top-Down DFS

### Intuition
A consecutive path grows only when a child's value is exactly `parent + 1`. So pass **down** into each call "the length of the increasing run so far, ending at me." If a child continues the run, its length is `parentLength + 1`; otherwise it resets to `1`. Track the global maximum along the way.

### Algorithm
1. `dfs(node, parentVal, lengthSoFar)`:
   - if `node` is nil, return.
   - `length = (node.Val == parentVal+1) ? lengthSoFar+1 : 1`.
   - update global `best` with `length`.
   - recurse left and right with `(node.Val, length)`.
2. Seed with `dfs(root, root.Val-1, 0)` so the root itself starts a run of length 1.

### Complexity
- **Time:** O(n) — each node visited exactly once.
- **Space:** O(h) — recursion stack proportional to tree height.

### Code
```go
func topDownDFS(root *TreeNode) int {
	best := 0
	var dfs func(node *TreeNode, parentVal, lengthSoFar int)
	dfs = func(node *TreeNode, parentVal, lengthSoFar int) {
		if node == nil {
			return
		}
		length := 1
		if node.Val == parentVal+1 {
			length = lengthSoFar + 1
		}
		if length > best {
			best = length
		}
		dfs(node.Left, node.Val, length)
		dfs(node.Right, node.Val, length)
	}
	if root != nil {
		dfs(root, root.Val-1, 0)
	}
	return best
}
```

### Dry Run
Example 1 tree: `1 → 3 → {2, 4 → 5}`. Seed `dfs(1, 0, 0)`.

| node | parentVal | node==parent+1? | length | best |
|------|-----------|-----------------|--------|------|
| 1    | 0         | no              | 1      | 1    |
| 3    | 1         | no (3≠2)        | 1      | 1    |
| 2    | 3         | no (2≠4)        | 1      | 1    |
| 4    | 3         | yes (4==3+1)    | 2      | 2    |
| 5    | 4         | yes (5==4+1)    | 3      | 3    |

Answer = **3** (path 3-4-5).

---

## Approach 2 — Bottom-Up DFS

### Intuition
Solve children first. The run **starting at** a node is `1 + child's run` **if** that child equals `node + 1`. Combine the left and right options, update the global best, and return the node's own run length upward.

### Algorithm
1. `dfs(node)` returns the length of the longest increasing run starting at `node`.
2. `len = 1`. If left exists and `left.Val == node.Val+1`, `len = max(len, 1+dfs(left))`. Same for right.
3. Update global `best` with `len`; return `len`.

### Complexity
- **Time:** O(n) — one visit per node.
- **Space:** O(h) — recursion depth.

### Code
```go
func bottomUpDFS(root *TreeNode) int {
	best := 0
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		leftLen := dfs(node.Left)
		rightLen := dfs(node.Right)

		length := 1
		if node.Left != nil && node.Left.Val == node.Val+1 {
			length = max(length, 1+leftLen)
		}
		if node.Right != nil && node.Right.Val == node.Val+1 {
			length = max(length, 1+rightLen)
		}
		if length > best {
			best = length
		}
		return length
	}
	dfs(root)
	return best
}
```

### Dry Run
Example 1, postorder (children before parents):

| node | leftLen | rightLen | consecutive child? | length returned | best |
|------|---------|----------|--------------------|-----------------|------|
| 2    | 0       | 0        | none               | 1               | 1    |
| 5    | 0       | 0        | none               | 1               | 1    |
| 4    | 0       | 1 (=5)   | right 5==4+1       | 1+1 = 2         | 2    |
| 3    | 1 (=2)  | 2 (=4)   | right 4==3+1       | 1+2 = 3         | 3    |
| 1    | 0       | 3 (=3)   | 3≠1+1              | 1               | 3    |

Answer = **3**.

---

## Key Takeaways
- **Consecutive-path problems on trees** reduce to a single DFS that either carries the running length down or returns it up.
- The reset rule is the crux: whenever `child != parent + 1`, the run restarts at `1`.
- Path direction is strictly top-down (parent→child), so you never combine a left run with a right run through the node — that would require going up then down (which is disallowed here, unlike the "any direction" variant #549).

---

## Related Problems
- LeetCode #549 — Binary Tree Longest Consecutive Sequence II (path may increase OR decrease, and may pass through a node child-to-child)
- LeetCode #128 — Longest Consecutive Sequence (array version)
- LeetCode #124 — Binary Tree Maximum Path Sum (bottom-up "return one side, combine both" pattern)
- LeetCode #687 — Longest Univalue Path
