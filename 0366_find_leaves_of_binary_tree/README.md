# 0366 — Find Leaves of Binary Tree

> LeetCode #366 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, collect a tree's nodes as if you were doing this:

- Collect all the leaf nodes.
- Remove all the leaf nodes.
- Repeat until the tree is empty.

**Example 1:**

```
Input: root = [1,2,3,4,5]
Output: [[4,5,3],[2],[1]]
Explanation:
[[3,5,4],[2],[1]] and [[3,4,5],[2],[1]] are also considered correct answers
since per each level it does not matter the order on which elements are returned.
```

**Example 2:**

```
Input: root = [1]
Output: [[1]]
```

**Constraints:**

- The number of nodes in the tree is in the range `[1, 100]`.
- `-100 <= Node.val <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| LinkedIn   | ★★★☆☆ Medium     | 2022          |
| Meta       | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Postorder DFS / Tree Traversal** — the pass in which a node is removed equals its height, and height is a bottom-up (postorder) quantity → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Graph DFS foundations** — recursion over children with values combined on the way back up → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Repeated Leaf Stripping (Brute Force) | O(n·h) | O(h) | Direct simulation; intuitive but re-walks the tree each pass |
| 2 | Postorder Height Grouping (Optimal) | O(n) | O(h) | The interview answer: one traversal, no mutation |

---

## Approach 1 — Repeated Leaf Stripping (Brute Force)

### Intuition

The statement is a procedure, so obey it. Each pass finds the current leaves, records their values as one group, and physically detaches them (nulls the parent→leaf pointers). Repeat until the root itself is removed. The root needs care because it has no parent — when the root becomes a leaf, the pass that collects it empties the tree.

### Algorithm

1. While the tree is non-empty:
   1. Run a DFS that, for every leaf encountered, appends its value to this pass's group and returns `nil` so the parent drops the pointer.
   2. The DFS returns the surviving (possibly `nil`) root back to the caller.
   3. Append this pass's collected group to the result.

### Complexity

- **Time:** O(n·h) — up to `h ≈ n` passes (a degenerate "stick" strips one leaf per pass), each pass walking O(n) nodes.
- **Space:** O(h) recursion stack per pass, plus O(n) output.

### Code

```go
func repeatedStripping(root *TreeNode) [][]int {
	var result [][]int // groups of values, one per stripping pass

	// removeLeaves detaches every current leaf below node, appending their
	// values to *collected. It returns the (possibly nil) node to keep — nil if
	// node itself was a leaf and should be removed by its caller.
	var removeLeaves func(node *TreeNode, collected *[]int) *TreeNode
	removeLeaves = func(node *TreeNode, collected *[]int) *TreeNode {
		if node == nil {
			return nil
		}
		if node.Left == nil && node.Right == nil {
			// node is a leaf: record it and tell the parent to drop it.
			*collected = append(*collected, node.Val)
			return nil
		}
		// Recurse first, then re-link the surviving children.
		node.Left = removeLeaves(node.Left, collected)
		node.Right = removeLeaves(node.Right, collected)
		return node
	}

	for root != nil {
		var pass []int // values collected in this single stripping pass
		root = removeLeaves(root, &pass)
		result = append(result, pass) // this pass forms one output group
	}
	return result
}
```

### Dry Run

Example 1: tree `[1,2,3,4,5]` — node 1 has children 2,3; node 2 has children 4,5.

| Pass | Current leaves | Group collected | Tree after pass |
|------|----------------|-----------------|-----------------|
| 1 | 4, 5, 3 | `[4,5,3]` | 1 → (2 → nil,nil), 3 removed |
| 2 | 2 | `[2]` | 1 with no children |
| 3 | 1 | `[1]` | empty |

Result: `[[4,5,3],[2],[1]]` ✔

---

## Approach 2 — Postorder Height Grouping (Optimal)

### Intuition

A node is removed on the pass equal to its **height** — the distance to its deepest leaf. Leaves (height 0) go first; a node whose deepest descendant leaf is one level down (height 1) goes second; and so on. Height is a classic postorder quantity: `height(node) = 1 + max(height(left), height(right))`, with `nil` treated as height −1 so a leaf computes 0. Use that height directly as the index of the output group. One traversal, no mutation, no repeated passes.

### Algorithm

1. DFS postorder. Return −1 for a `nil` child.
2. `h = 1 + max(dfs(left), dfs(right))`.
3. If `h == len(result)`, this is the first node seen at that height — open a new group.
4. Append `node.Val` to `result[h]`; return `h`.

### Complexity

- **Time:** O(n) — each node visited exactly once.
- **Space:** O(h) recursion stack (O(n) worst case) plus O(n) output.

### Code

```go
func heightGrouping(root *TreeNode) [][]int {
	var result [][]int // result[h] holds every node of height h

	// dfs returns the height of node (leaf = 0, nil = -1) and files node into
	// the group indexed by its height along the way.
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return -1 // so a leaf computes 1 + max(-1,-1) = 0
		}
		left := dfs(node.Left)              // height of left subtree
		right := dfs(node.Right)            // height of right subtree
		h := 1 + max(left, right)           // this node's height
		if h == len(result) {               // first node discovered at this height
			result = append(result, []int{}) // open a new group
		}
		result[h] = append(result[h], node.Val) // file node under its height
		return h
	}

	dfs(root)
	return result
}
```

### Dry Run

Example 1: tree `[1,2,3,4,5]`. Postorder visits 4, 5, 2, 3, 1.

| Visit | node | dfs(left) | dfs(right) | h = 1+max | result after |
|-------|------|-----------|------------|-----------|--------------|
| 1 | 4 | -1 | -1 | 0 | `[[4]]` |
| 2 | 5 | -1 | -1 | 0 | `[[4,5]]` |
| 3 | 2 | 0 (from 4) | 0 (from 5) | 1 | `[[4,5],[2]]` |
| 4 | 3 | -1 | -1 | 0 | `[[4,5,3],[2]]` |
| 5 | 1 | 1 (from 2) | 0 (from 3) | 2 | `[[4,5,3],[2],[1]]` |

Result: `[[4,5,3],[2],[1]]` ✔

---

## Key Takeaways

- **Removal order = node height.** Whenever a problem repeatedly peels "outermost" elements, ask what monotone quantity indexes the peel — here it is height, computable in one postorder pass.
- Use `nil → -1` as the DFS base case so leaves naturally land at height 0.
- Indexing `result[h]` while growing the slice on first sighting (`h == len(result)`) is a clean way to bucket-by-depth without pre-sizing.
- The brute-force simulation is worth knowing to explain *why* the height trick is correct — they compute the same groups.

---

## Related Problems

- LeetCode #104 — Maximum Depth of Binary Tree (the same height recurrence)
- LeetCode #110 — Balanced Binary Tree (heights compared bottom-up)
- LeetCode #543 — Diameter of Binary Tree (postorder combining child heights)
- LeetCode #669 — Trim a Binary Tree (returning a rewired subtree from DFS)
