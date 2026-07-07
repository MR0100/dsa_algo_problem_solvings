# 0427 — Construct Quad Tree

> LeetCode #427 · Difficulty: Medium
> **Categories:** Array, Divide and Conquer, Tree, Matrix

---

## Problem Statement

Given a `n * n` matrix `grid` of `0`'s and `1`'s only. We want to represent `grid` with a Quad-Tree.

Return *the root of the Quad-Tree representing* `grid`.

A Quad-Tree is a tree data structure in which each internal node has exactly four children. Besides, each node has two attributes:

- `val`: **True** if the node represents a grid of 1's or **False** if the node represents a grid of 0's. Notice that you can assign the `val` to True or False when `isLeaf` is **False**, and both are accepted in the answer.
- `isLeaf`: **True** if the node is a leaf node on the tree or **False** if the node has four children.

```
class Node {
    public boolean val;
    public boolean isLeaf;
    public Node topLeft;
    public Node topRight;
    public Node bottomLeft;
    public Node bottomRight;
}
```

We can construct a Quad-Tree from a two-dimensional area using the following steps:

1. If the current grid has the same value (i.e all `1`'s or all `0`'s) set `isLeaf` True and set `val` to the value of the grid and set the four children to Null and stop.
2. If the current grid has different values, set `isLeaf` to False and set `val` to any value and divide the current grid into four sub-grids as shown in the photo.
3. Recurse for each of the children with the proper sub-grid.

If you want to know more about the Quad-Tree, you can refer to the [wiki](https://en.wikipedia.org/wiki/Quadtree).

**Quad-Tree format:** The output represents the serialized format of a Quad-Tree using level order traversal, where `null` signifies a path terminator where no node exists below.

It is very similar to the serialization of the binary tree. The only difference is that the node is represented as a list `[isLeaf, val]`.

If the value of `isLeaf` or `val` is True we represent it as **1** in the list `[isLeaf, val]` and if the value of `isLeaf` or `val` is False we represent it as **0**.

**Example 1:**
```
Input: grid = [[0,1],[1,0]]
Output: [[0,1],[1,0],[1,1],[1,1],[1,0]]
```
Explanation: The explanation of this example is shown below:
Notice that 0 represents False and 1 represents True in the photo representing the Quad-Tree. The topLeft, topRight, bottomLeft and bottomRight cells each differ, so the root splits into four 1×1 leaves.

**Example 2:**
```
Input: grid = [[1,1,1,1,0,0,0,0],
               [1,1,1,1,0,0,0,0],
               [1,1,1,1,1,1,1,1],
               [1,1,1,1,1,1,1,1],
               [1,1,1,1,0,0,0,0],
               [1,1,1,1,0,0,0,0],
               [1,1,1,1,0,0,0,0],
               [1,1,1,1,0,0,0,0]]
Output: [[0,1],[1,1],[0,1],[1,1],[1,0],null,null,null,null,[1,0],[1,0],[1,1],[1,1]]
```
Explanation: All values in the grid are not the same. We divide the grid into four sub-grids. The topLeft, bottomLeft and bottomRight each has the same value. The topRight have different values so we divide it into 4 sub-grids so that each sub-grid has the same value.

**Constraints:**
- `n == grid.length == grid[i].length`
- `n == 2ˣ` where `0 <= x <= 6`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Divide and Conquer** — the defining pattern: split each square into four equal quadrants, solve each, combine → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Matrix / 2-D grid indexing** — quadrants are addressed by `(row, col, size)` offsets into the same grid, no copying → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Tree construction (level-order serialization)** — the result is a 4-ary tree; verifying it uses BFS serialization → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)

---

## Approaches Overview

Let n = side length of the grid (a power of two).

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Divide & conquer, scan each sub-grid for uniformity | O(n² log n) | O(log n) stack + O(#nodes) | Simplest to reason about; fine for the tiny n ≤ 64 here |
| 2 | Divide & conquer, merge four leaf children (Optimal) | O(n²) | O(log n) stack + O(#nodes) | Reads each cell once; the clean linear-work solution |

---

## Approach 1 — Divide & Conquer (Scan Each Sub-grid)

### Intuition
A square region can be one leaf **iff every cell in it is identical**. So test that directly: if the current `size × size` block is uniform, emit a leaf carrying that value. If not, the region truly needs four children — split into the four equal quadrants (each of side `size/2`, always possible because `n` is a power of two) and recurse. Wrap the four results in an internal node. The recursion bottoms out at `1 × 1` blocks, which are trivially uniform.

### Algorithm
1. `build(r, c, size)`: if the block at `(r,c)` of side `size` is uniform, return a leaf with `Val = grid[r][c]==1`.
2. Otherwise `half = size/2`; recurse on the four quadrants: TopLeft `(r,c)`, TopRight `(r,c+half)`, BottomLeft `(r+half,c)`, BottomRight `(r+half,c+half)`.
3. Return an internal node (`IsLeaf=false`) holding the four children.

### Complexity
- **Time:** O(n² log n) worst case — a fully-mixed grid re-scans the region at every level; the log n levels each do O(n²) total scanning.
- **Space:** O(log n) recursion depth plus O(#nodes) for the constructed tree.

### Code
```go
func divideAndConquer(grid [][]int) *Node {
	n := len(grid)
	var build func(r, c, size int) *Node
	build = func(r, c, size int) *Node {
		if isUniform(grid, r, c, size) {
			// Whole block equal ⇒ one leaf. grid[r][c]==1 → Val true.
			return &Node{Val: grid[r][c] == 1, IsLeaf: true}
		}
		half := size / 2 // each quadrant is half the side length
		return &Node{
			Val:         true, // arbitrary for internal nodes; true is conventional
			IsLeaf:      false,
			TopLeft:     build(r, c, half),
			TopRight:    build(r, c+half, half),
			BottomLeft:  build(r+half, c, half),
			BottomRight: build(r+half, c+half, half),
		}
	}
	return build(0, 0, n)
}

// isUniform reports whether every cell in the size×size block at (r,c) equals
// the top-left cell of that block.
func isUniform(grid [][]int, r, c, size int) bool {
	first := grid[r][c] // reference value for the block
	for i := r; i < r+size; i++ {
		for j := c; j < c+size; j++ {
			if grid[i][j] != first {
				return false // found a differing cell ⇒ not uniform
			}
		}
	}
	return true
}
```

### Dry Run (Example 1)

`grid = [[0,1],[1,0]]`, call `build(0,0,2)`.

| Call | Block cells | Uniform? | Result |
|------|-------------|----------|--------|
| `build(0,0,2)` | 0,1,1,0 | no | split into 4, internal `[0,1]` |
| `build(0,0,1)` TopLeft | 0 | yes | leaf `[1,0]` |
| `build(0,1,1)` TopRight | 1 | yes | leaf `[1,1]` |
| `build(1,0,1)` BottomLeft | 1 | yes | leaf `[1,1]` |
| `build(1,1,1)` BottomRight | 0 | yes | leaf `[1,0]` |

Level-order serialization (leaf children are trailing nulls, trimmed): `[[0,1],[1,0],[1,1],[1,1],[1,0]]` ✓

---

## Approach 2 — Divide & Conquer Merging Children (Optimal)

### Intuition
Scanning a whole block to test uniformity (Approach 1) repeats work: an outer mixed block scans cells that its children will scan again. Flip it around — recurse **first**, all the way down to `1 × 1` leaves, then ask a purely local question at each internal step: *"are my four children all leaves with the same value?"* If yes, the block was uniform after all, so collapse them into a single merged leaf; if no, keep them. Every cell is read exactly once (at its base case), so total work is linear in the grid.

### Algorithm
1. `build(r, c, size)`: if `size == 1`, return a leaf for that single cell.
2. Recurse into the four `size/2` quadrants → `tl, tr, bl, br`.
3. If all four are leaves **and** `tl.Val == tr.Val == bl.Val == br.Val`, return one merged leaf with that value (discard the children).
4. Otherwise return an internal node holding the four children.

### Complexity
- **Time:** O(n²) — each cell contributes to exactly one base-case leaf; every internal node does O(1) merge work, and there are O(n²) nodes total.
- **Space:** O(log n) recursion depth plus O(#nodes) for the tree.

### Code
```go
func mergeChildren(grid [][]int) *Node {
	var build func(r, c, size int) *Node
	build = func(r, c, size int) *Node {
		if size == 1 {
			// Base case: a single cell is always a leaf.
			return &Node{Val: grid[r][c] == 1, IsLeaf: true}
		}
		half := size / 2
		tl := build(r, c, half)
		tr := build(r, c+half, half)
		bl := build(r+half, c, half)
		br := build(r+half, c+half, half)

		// Collapse iff all four children are leaves with an identical value.
		if tl.IsLeaf && tr.IsLeaf && bl.IsLeaf && br.IsLeaf &&
			tl.Val == tr.Val && tr.Val == bl.Val && bl.Val == br.Val {
			return &Node{Val: tl.Val, IsLeaf: true} // merged leaf; children discarded
		}
		return &Node{ // genuinely mixed region: keep the four children
			Val:         true,
			IsLeaf:      false,
			TopLeft:     tl,
			TopRight:    tr,
			BottomLeft:  bl,
			BottomRight: br,
		}
	}
	return build(0, 0, len(grid))
}
```

### Dry Run (Example 1)

`grid = [[0,1],[1,0]]`, call `build(0,0,2)`.

| Call | size | Recurse / decide | Result |
|------|------|------------------|--------|
| `build(0,0,1)` | 1 | base case, cell 0 | leaf `Val=false` `[1,0]` |
| `build(0,1,1)` | 1 | base case, cell 1 | leaf `Val=true` `[1,1]` |
| `build(1,0,1)` | 1 | base case, cell 1 | leaf `Val=true` `[1,1]` |
| `build(1,1,1)` | 1 | base case, cell 0 | leaf `Val=false` `[1,0]` |
| `build(0,0,2)` | 2 | children all leaves but Vals {F,T,T,F} differ → **no merge** | internal `[0,1]` with those 4 children |

Serialization: `[[0,1],[1,0],[1,1],[1,1],[1,0]]` ✓ — identical tree to Approach 1, built bottom-up.

---

## Key Takeaways

- **Two dual framings of the same recursion:** top-down "is this block uniform? if not, split" (Approach 1) vs bottom-up "split to leaves, then merge equal siblings" (Approach 2). The merge version reads each cell once and drops the extra `log n` factor.
- **Powers of two make the split clean** — `size/2` always partitions into four equal squares down to `1×1`; no bounds fiddling or ragged quadrants.
- **Address quadrants by offset, never copy the sub-grid** — pass `(r, c, size)` into the same matrix. Copying would blow the memory and time up by `O(n²)` per level.
- **`val` is irrelevant for internal nodes** — the checker only reads `val` when `isLeaf` is true, so any value is accepted there (this repo uses `true`).
- The LeetCode serializer gives every node four child slots (leaves' are `null`) and trims trailing `null`s — worth replicating exactly when you self-verify, since a "compact" BFS produces a different (also-correct-looking) string that won't match the expected output.

---

## Related Problems
- LeetCode #558 — Logical OR of Two Binary Grids Represented as Quad-Trees (operates on this exact structure)
- LeetCode #427 companion #395/#241 — divide-and-conquer decomposition on other inputs
- LeetCode #240 — Search a 2D Matrix II (quadrant-style pruning of a grid)
- LeetCode #308 — Range Sum Query 2D (hierarchical decomposition of a grid)
- LeetCode #105 — Construct Binary Tree from Preorder and Inorder (recursive tree construction)
