# 0314 — Binary Tree Vertical Order Traversal

> LeetCode #314 · Difficulty: Medium
> **Categories:** Hash Table, Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return the **vertical order traversal** of its nodes' values. (i.e., from top to bottom, column by column).

If two nodes are in the same row and column, the order should be from **left to right**.

**Example 1:**

```
Input: root = [3,9,20,null,null,15,7]
Output: [[9],[3,15],[20],[7]]
```

**Example 2:**

```
Input: root = [3,9,8,4,0,1,7]
Output: [[4],[9],[3,0,1],[8],[7]]
```

**Example 3:**

```
Input: root = [3,9,8,4,0,1,7,null,null,null,2,5]
Output: [[4],[9,5],[3,0,1],[8,2],[7]]
```

**Constraints:**

- The number of nodes in the tree is in the range `[0, 100]`.
- `-100 <= Node.val <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Meta       | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **BFS (level-order traversal)** — visiting nodes top-to-bottom, left-before-right, which matches the required within-column order for free → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Tree traversal** — carrying a column index down through children → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Hash Map** — mapping column index → list of node values → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Queue / Deque** — FIFO queue driving the BFS → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS with Column Index (Optimal) | O(n) | O(n) | Standard answer; ordering is automatic |
| 2 | DFS with (col,row) + stable sort | O(n log n) | O(n) | When a recursive walk is preferred |

---

## Approach 1 — BFS with Column Index (Optimal)

### Intuition
Vertical order groups nodes by **column**: root = 0, left child = `col-1`, right child = `col+1`. Within a column, nodes must be top-to-bottom and left-before-right on ties. A BFS visits nodes strictly top-to-bottom and, by enqueuing the left child **before** the right child, left-before-right within each level — exactly the tie order required. So appending each dequeued node to its column bucket gives the correct order with **no sorting**.

### Algorithm
1. BFS from root, carrying each node's column in the queue.
2. Track `minCol`/`maxCol`; append each value to `cols[col]` in visit order.
3. Enqueue left child before right child.
4. Emit columns from `minCol` to `maxCol`.

### Complexity
- **Time:** O(n) — each node visited once; O(range) to assemble output.
- **Space:** O(n) — the queue and the column buckets.

### Code
```go
func bfs(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	cols := map[int][]int{}
	minCol, maxCol := 0, 0

	type item struct {
		node *TreeNode
		col  int
	}
	queue := []item{{root, 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		cols[cur.col] = append(cols[cur.col], cur.node.Val)
		if cur.col < minCol {
			minCol = cur.col
		}
		if cur.col > maxCol {
			maxCol = cur.col
		}

		if cur.node.Left != nil {
			queue = append(queue, item{cur.node.Left, cur.col - 1})
		}
		if cur.node.Right != nil {
			queue = append(queue, item{cur.node.Right, cur.col + 1})
		}
	}

	res := make([][]int, 0, maxCol-minCol+1)
	for c := minCol; c <= maxCol; c++ {
		res = append(res, cols[c])
	}
	return res
}
```

### Dry Run
`root = [3,9,20,null,null,15,7]`. Queue processed FIFO; column shown per node.

| dequeue | col | cols after | minCol,maxCol | enqueue |
|---------|-----|------------|---------------|---------|
| 3 | 0 | {0:[3]} | 0,0 | 9@-1, 20@+1 |
| 9 | -1 | {-1:[9], 0:[3]} | -1,0 | (leaf) |
| 20 | +1 | {..., 1:[20]} | -1,1 | 15@0, 7@+2 |
| 15 | 0 | {0:[3,15]} | -1,1 | (leaf) |
| 7 | +2 | {2:[7]} | -1,2 | (leaf) |

Emit cols -1..2: `[[9],[3,15],[20],[7]]`.

---

## Approach 2 — DFS with (col, row) + Stable Sort

### Intuition
DFS is a natural recursive walk but dives deep before wide, so it does **not** visit strictly top-to-bottom. Record each node's `(column, row)` where row = depth. To rebuild vertical order, sort each column's entries by row. A **stable** sort keeps left-before-right for equal `(col,row)` because we recurse left before right, appending left nodes first.

### Algorithm
1. DFS carrying `(col, row)`; append `(row, val)` into `cols[col]`.
2. For each column, stable-sort by row.
3. Emit columns from `minCol` to `maxCol`.

### Complexity
- **Time:** O(n log n) — dominated by per-column sorting.
- **Space:** O(n) — stored entries and recursion stack.

### Code
```go
func dfsSort(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	type rowVal struct {
		row int
		val int
	}
	cols := map[int][]rowVal{}
	minCol, maxCol := 0, 0

	var dfs func(node *TreeNode, col, row int)
	dfs = func(node *TreeNode, col, row int) {
		if node == nil {
			return
		}
		cols[col] = append(cols[col], rowVal{row, node.Val})
		if col < minCol {
			minCol = col
		}
		if col > maxCol {
			maxCol = col
		}
		dfs(node.Left, col-1, row+1)
		dfs(node.Right, col+1, row+1)
	}
	dfs(root, 0, 0)

	res := make([][]int, 0, maxCol-minCol+1)
	for c := minCol; c <= maxCol; c++ {
		entries := cols[c]
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].row < entries[j].row
		})
		vals := make([]int, len(entries))
		for i, e := range entries {
			vals[i] = e.val
		}
		res = append(res, vals)
	}
	return res
}
```

### Dry Run
`root = [3,9,20,null,null,15,7]`. DFS order: 3, 9, 20, 15, 7.

| visit | (col,row) | cols after append |
|-------|-----------|-------------------|
| 3 | (0,0) | {0:[(0,3)]} |
| 9 | (-1,1) | {-1:[(1,9)]} |
| 20 | (1,1) | {1:[(1,20)]} |
| 15 | (0,2) | {0:[(0,3),(2,15)]} |
| 7 | (2,2) | {2:[(2,7)]} |

Per-column stable sort by row: col 0 → `[3,15]` (rows 0,2). Emit cols -1..2: `[[9],[3,15],[20],[7]]`.

---

## Key Takeaways
- **BFS gives ordering for free.** Because BFS is top-down and (with left-before-right enqueue) left-first, no row-based sorting is needed — a key advantage over DFS for this problem.
- **Column index = root 0, left −1, right +1**, carried through the traversal, is the core idea of all vertical/diagonal tree problems.
- **DFS needs an explicit row** and a stable sort to recover top-to-bottom order, since it is not level-ordered.
- Contrast with LeetCode #987 (Vertical Order Traversal of a Binary Tree), where same `(row,col)` nodes are additionally sorted by **value** — here they are sorted by left-to-right position instead.

---

## Related Problems
- LeetCode #987 — Vertical Order Traversal of a Binary Tree (adds value tie-break)
- LeetCode #102 — Binary Tree Level Order Traversal (BFS with buckets)
- LeetCode #199 — Binary Tree Right Side View (BFS per level)
- LeetCode #103 — Binary Tree Zigzag Level Order Traversal (BFS ordering)
