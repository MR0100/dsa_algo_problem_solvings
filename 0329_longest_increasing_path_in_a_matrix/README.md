# 0329 — Longest Increasing Path in a Matrix

> LeetCode #329 · Difficulty: Hard
> **Categories:** Dynamic Programming, Depth-First Search, Breadth-First Search, Graph, Topological Sort, Memoization, Matrix

---

## Problem Statement

Given an `m x n` integers `matrix`, return *the length of the longest increasing path in* `matrix`.

From each cell, you can either move in four directions: left, right, up, or down. You **may not** move **diagonally** or move **outside the boundary** (i.e., wrap-around is not allowed).

**Example 1:**

```
Input: matrix = [[9,9,4],[6,6,8],[2,1,1]]
Output: 4
Explanation: The longest increasing path is [1, 2, 6, 9].
```

**Example 2:**

```
Input: matrix = [[3,4,5],[3,2,6],[2,2,1]]
Output: 4
Explanation: The longest increasing path is [3, 4, 5, 6]. Moving diagonally is not allowed.
```

**Example 3:**

```
Input: matrix = [[1]]
Output: 1
```

**Constraints:**

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 200`
- `0 <= matrix[i][j] <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **DFS on a grid** — the natural recursion is "longest increasing path starting here = 1 + best over strictly-larger neighbours"; a depth-first walk explores each such continuation → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **2D Dynamic Programming / Memoization** — because values strictly increase along any path, the cell-to-larger-neighbour graph is a DAG, so `memo[i][j]` (the answer for a cell) is independent of how you reached it and can be cached, collapsing exponential DFS to O(m·n) → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Topological Sort (Kahn's BFS)** — treating the grid as a DAG, peeling off out-degree-0 "peak" cells layer by layer counts the longest chain without recursion → see [`/dsa/topological_sort.md`](/dsa/topological_sort.md)
- **Matrix Traversal** — four-directional neighbour scanning with a `{{-1,0},{1,0},{0,-1},{0,1}}` direction array and boundary checks underpins every approach → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Plain DFS (Brute Force) | O(2^(m·n)) | O(m·n) | Baseline / intuition only; TLEs on large grids because overlapping subpaths are recomputed |
| 2 | DFS + Memoization (Optimal) | O(m·n) | O(m·n) | The standard interview answer; shortest to write and prove |
| 3 | Topological Sort (BFS Peeling) | O(m·n) | O(m·n) | Iterative alternative with no recursion depth risk; showcases the DAG framing |

---

## Approach 1 — Plain DFS (Brute Force)

### Intuition

The longest increasing path starting at a cell `(i,j)` is `1` plus the longest increasing path starting at whichever strictly-larger neighbour extends best. That recursion is directly executable: from every cell, DFS into all four neighbours that hold a strictly larger value and keep the best `1 + child`. Since edges only ever go small→large, a value can never repeat inside one path, so we need **no visited set**. The catch is that with no caching, the same cell's subtree is re-explored every time a path passes through it, making the work exponential — correct, but a TLE on big inputs.

### Algorithm

1. For each cell `(i,j)`, compute `dfs(i,j)` = length of the longest increasing path that **starts** at `(i,j)`.
2. Inside `dfs`, start `best = 1` (the cell itself). For each of the four neighbours that is in bounds and strictly larger, recurse and take `1 + dfs(neighbour)`, keeping the maximum.
3. Return `best`.
4. The answer is the maximum `dfs(i,j)` over all cells.

### Complexity

- **Time:** O(2^(m·n)) worst case — with no memo, overlapping increasing subpaths are recomputed from scratch on every visit.
- **Space:** O(m·n) — the recursion stack can be as deep as the longest path, which is at most `m·n` cells.

### Code

```go
func dfsBrute(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])

	// dfs returns the length of the longest increasing path that STARTS at (i,j).
	var dfs func(i, j int) int
	dfs = func(i, j int) int {
		best := 1 // the cell itself is a path of length 1
		for _, d := range directions {
			ni, nj := i+d[0], j+d[1] // neighbour coordinates
			// Must stay in bounds AND be strictly larger to extend the path.
			if ni >= 0 && ni < m && nj >= 0 && nj < n && matrix[ni][nj] > matrix[i][j] {
				length := 1 + dfs(ni, nj) // this cell + best path from neighbour
				if length > best {
					best = length
				}
			}
		}
		return best
	}

	ans := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if v := dfs(i, j); v > ans { // try starting from every cell
				ans = v
			}
		}
	}
	return ans
}
```

### Dry Run

Example 1: `matrix = [[9,9,4],[6,6,8],[2,1,1]]`. Cells are `(row,col)`; the winning path is `1 → 2 → 6 → 9`, i.e. `(2,1) → (2,0) → (1,0) → (0,0)`.

Trace `dfs(2,1)` (the cell holding value `1`):

| Call | Cell | Value | Strictly-larger neighbours explored | Returns |
|------|------|-------|-------------------------------------|---------|
| `dfs(2,1)` | (2,1) | 1 | up (1,1)=6 ✔; left (2,0)=2 ✔; right (2,2)=1 ✗ | 1 + max(child paths) = **4** |
| ↳ `dfs(2,0)` | (2,0) | 2 | up (1,0)=6 ✔ | 1 + 3 = **3** |
| ↳ `dfs(1,0)` | (1,0) | 6 | up (0,0)=9 ✔; right (1,1)=6 ✗ | 1 + 2 = **2** |
| ↳ `dfs(0,0)` | (0,0) | 9 | none larger | **1** |
| ↳ `dfs(1,1)` | (1,1) | 6 | up (0,1)=9 ✔; right (1,2)=8 ✔ | 1 + 2 = **3** |

`dfs(2,1)` compares `1 + dfs(1,1)=4` vs `1 + dfs(2,0)=4` → returns `4`. Scanning all starting cells, the maximum found is `4`, produced by the path `1 → 2 → 6 → 9`. ✔ Note `dfs(0,0)`, `dfs(1,0)` etc. are recomputed whenever another path reaches them — that redundancy is exactly what Approach 2 removes.

---

## Approach 2 — DFS + Memoization (Optimal)

### Intuition

Draw an edge from every cell to each strictly-larger orthogonal neighbour. Because values **strictly** increase along any edge, no path can ever return to a cell — the graph is a Directed Acyclic Graph (DAG). In a DAG, the longest path starting at a node depends only on its successors, never on how you arrived, so the quantity `memo[i][j] = longest increasing path starting at (i,j)` is well-defined and safe to cache. Run the same DFS as Approach 1, but the first time a cell is solved, store the result; every later visit reads the cache in O(1). Each cell is computed exactly once, turning the exponential blow-up into linear work.

### Algorithm

1. Allocate `memo` initialised to `0`, where `0` means "not computed yet".
2. `dfs(i,j)`: if `memo[i][j] != 0`, return it. Otherwise set `best = 1` and, for each in-bounds strictly-larger neighbour, take `1 + dfs(neighbour)`, keeping the max.
3. Store `best` in `memo[i][j]` before returning.
4. Return the maximum `dfs(i,j)` over all cells.

### Complexity

- **Time:** O(m·n) — each cell is fully computed once; every cell does O(1) work (four neighbour checks), and cached cells return immediately.
- **Space:** O(m·n) — the `memo` table plus recursion stack bounded by the longest path.

### Code

```go
func dfsMemo(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])

	// memo[i][j] = longest increasing path starting at (i,j); 0 = uncomputed.
	memo := make([][]int, m)
	for i := range memo {
		memo[i] = make([]int, n)
	}

	var dfs func(i, j int) int
	dfs = func(i, j int) int {
		if memo[i][j] != 0 {
			return memo[i][j] // already solved this cell — reuse it
		}
		best := 1 // the cell alone
		for _, d := range directions {
			ni, nj := i+d[0], j+d[1]
			if ni >= 0 && ni < m && nj >= 0 && nj < n && matrix[ni][nj] > matrix[i][j] {
				length := 1 + dfs(ni, nj)
				if length > best {
					best = length
				}
			}
		}
		memo[i][j] = best // cache before returning so callers reuse it
		return best
	}

	ans := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if v := dfs(i, j); v > ans {
				ans = v
			}
		}
	}
	return ans
}
```

### Dry Run

Example 1: `matrix = [[9,9,4],[6,6,8],[2,1,1]]`. We scan cells in row-major order `(0,0), (0,1), …`; each `dfs` fills `memo` the first time it resolves a cell. Below is the fill order for the cells on and around the winning path.

| Order filled | Cell (val) | Larger neighbours | `memo[i][j]` | Reason |
|--------------|------------|-------------------|--------------|--------|
| 1 | (0,0)=9 | none | **1** | peak, dead end |
| 2 | (0,1)=9 | none | **1** | peak |
| 3 | (0,2)=4 | (1,2)=8 | **2** | 1 + memo(1,2)… |
| 4 | (1,2)=8 | (0,2)? 4<8 no; none larger | **1** | 8 is a peak locally |
| — | (0,2) resolves | uses memo(1,2)=1 | **2** | 4 → 8 |
| 5 | (1,0)=6 | (0,0)=9 → memo 1 | **2** | 6 → 9 |
| 6 | (1,1)=6 | (0,1)=9→1; (1,2)=8→1 | **2** | 6 → 9 (or 6 → 8) |
| 7 | (2,0)=2 | (1,0)=6 → memo 2 | **3** | 2 → 6 → 9, reuses cached (1,0) |
| 8 | (2,1)=1 | (1,1)=6→2; (2,0)=2→3 | **4** | 1 + max(2,3) = 4 |
| 9 | (2,2)=1 | (1,2)=8→1 | **2** | 1 → 8 |

When `dfs(2,1)` needs `dfs(2,0)` and `dfs(1,0)`, those are already cached (`3` and `2`), so no recomputation happens — this is the whole win over Approach 1. The maximum `memo` value is `4` at `(2,1)`, matching the path `1 → 2 → 6 → 9`. ✔

---

## Approach 3 — Topological Sort (BFS Peeling)

### Intuition

Keep the DAG framing but solve it iteratively with Kahn's algorithm. Orient every edge small→large; a cell's **out-degree** is how many neighbours are strictly larger. A cell with out-degree `0` is a local **peak** — no increasing path continues past it, so it must be the *last* cell of some increasing path. Remove all current peaks at once (one BFS layer); doing so deletes their incoming edges, which may drop some predecessor's out-degree to `0`, exposing the next layer of peaks. Each layer you peel corresponds to advancing one step deeper along the longest chain, so the number of layers equals the length of the longest increasing path.

### Algorithm

1. Compute `outDegree[i][j]` = number of strictly-larger neighbours of `(i,j)`.
2. Enqueue every cell whose `outDegree` is `0` (the initial peaks).
3. Process the queue layer by layer. For each popped cell, look at its **smaller** neighbours (its predecessors), decrement their `outDegree`, and enqueue any that reach `0` for the next layer.
4. Count the layers; that count is the answer.

### Complexity

- **Time:** O(m·n) — computing out-degrees scans four neighbours per cell, and each cell is enqueued and processed exactly once.
- **Space:** O(m·n) — the `outDegree` grid and the BFS queue.

### Code

```go
func topoSortBFS(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])

	// outDegree[i][j] = number of strictly-larger orthogonal neighbours.
	outDegree := make([][]int, m)
	for i := range outDegree {
		outDegree[i] = make([]int, n)
	}

	queue := make([][2]int, 0, m*n) // holds coordinates of current-layer peaks
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			for _, d := range directions {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < m && nj >= 0 && nj < n && matrix[ni][nj] > matrix[i][j] {
					outDegree[i][j]++ // (i,j) has an outgoing edge to a bigger neighbour
				}
			}
			if outDegree[i][j] == 0 {
				queue = append(queue, [2]int{i, j}) // a peak: nowhere larger to go
			}
		}
	}

	layers := 0 // how many BFS layers we peel = length of the longest path
	for len(queue) > 0 {
		layers++
		next := make([][2]int, 0) // cells that become peaks after this layer
		for _, cell := range queue {
			i, j := cell[0], cell[1]
			for _, d := range directions {
				ni, nj := i+d[0], j+d[1]
				// Look at SMALLER neighbours — they point INTO this cell.
				if ni >= 0 && ni < m && nj >= 0 && nj < n && matrix[ni][nj] < matrix[i][j] {
					outDegree[ni][nj]-- // remove the edge into the just-peeled cell
					if outDegree[ni][nj] == 0 {
						next = append(next, [2]int{ni, nj}) // newly exposed peak
					}
				}
			}
		}
		queue = next // advance to the next layer
	}
	return layers
}
```

### Dry Run

Example 1: `matrix = [[9,9,4],[6,6,8],[2,1,1]]`.

Initial out-degrees (count of strictly-larger neighbours):

```
val:  9 9 4        outDeg:  0 0 1
      6 6 8                 1 1 0
      2 1 1                 1 2 1
```

- `(0,0)=9`, `(0,1)=9`, `(1,2)=8` have out-degree `0` → **Layer 1** peaks.

| Layer | Peeled cells (out-degree 0) | Effect on predecessors' out-degree |
|-------|-----------------------------|-------------------------------------|
| 1 | (0,0)=9, (0,1)=9, (1,2)=8 | smaller neighbours: (0,2)=4 → 0, (1,0)=6 → 0, (1,1)=6 → 0, (2,2)=1 → 0 |
| 2 | (0,2)=4, (1,0)=6, (1,1)=6, (2,2)=1 | smaller neighbours: (2,0)=2 → 0, (2,1)=1 → 1 |
| 3 | (2,0)=2 | smaller neighbour (2,1)=1 → 0 |
| 4 | (2,1)=1 | no smaller neighbours left |

Four layers were peeled, so the answer is `4`. The layer index at which each cell is removed equals the length of the longest path *ending* at it: `(2,1)=1` leaves in layer 4 — exactly the chain `1 → 2 → 6 → 9`. ✔

---

## Key Takeaways

- **"Strictly increasing" ⇒ DAG.** Whenever the moves you may take are governed by a strict inequality on cell values, the induced graph is acyclic, which is the license to memoize (a cell's answer can't depend on the path taken to reach it).
- **Longest path in a DAG has two dual solutions:** top-down DFS + memo (Approach 2) and bottom-up topological peeling (Approach 3). Both are O(V + E) = O(m·n) here since each cell has ≤ 4 edges.
- **No visited set needed** for the DFS: the strict-increase rule already forbids revisiting a cell within a single path — a common point of confusion versus flood-fill problems where you *do* track visited.
- **Memoization is the single change** that takes this from exponential to linear. The recursion is identical to brute force; only the cache differs. Recognising the overlapping-subproblem structure is the whole interview.
- **Kahn's out-degree peeling** counts the longest chain as "number of layers" — a clean, recursion-free technique worth having when stack depth (up to `m·n = 40000`) is a concern.

---

## Related Problems

- LeetCode #200 — Number of Islands (grid DFS/BFS traversal with four-directional moves)
- LeetCode #62 — Unique Paths (2D DP over a grid)
- LeetCode #547 — Number of Provinces (graph connectivity / DFS)
- LeetCode #417 — Pacific Atlantic Water Flow (grid DFS with monotone-value flow, same neighbour-comparison pattern)
