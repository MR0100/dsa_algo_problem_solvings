# 0463 — Island Perimeter

> LeetCode #463 · Difficulty: Easy
> **Categories:** Array, Matrix, Depth-First Search, Breadth-First Search

---

## Problem Statement

You are given `row x col` `grid` representing a map where `grid[i][j] = 1` represents land and `grid[i][j] = 0` represents water.

Grid cells are connected **horizontally/vertically** (not diagonally). The `grid` is completely surrounded by water, and there is exactly **one** island (i.e., one or more connected land cells).

The island doesn't have "lakes", meaning the water inside isn't connected to the water around the island. One cell is a square with side length `1`. The grid is rectangular, width and height don't exceed `100`. Determine the perimeter of the island.

**Example 1:**

```
Input: grid = [[0,1,0,0],[1,1,1,0],[0,1,0,0],[1,1,0,0]]
Output: 16
Explanation: The perimeter is the 16 yellow stripes in the image below.
```

**Example 2:**

```
Input: grid = [[1]]
Output: 4
```

**Example 3:**

```
Input: grid = [[1,0]]
Output: 4
```

**Constraints:**

- `row == grid.length`
- `col == grid[i].length`
- `1 <= row, col <= 100`
- `grid[i][j]` is `0` or `1`.
- There is exactly one island in `grid`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Matrix Traversal** — the core operation is scanning a 2-D grid and inspecting each cell's four orthogonal neighbours with bounds checks → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Graph BFS/DFS** — the island is a connected component; a flood fill over land cells, counting sides that face water/edge, is the graph-search framing → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Cell Scan (count exposed sides) | O(rows · cols) | O(1) | Most direct; check 4 neighbours per land cell |
| 2 | Contribution Counting (Optimal) | O(rows · cols) | O(1) | Cleanest one-pass; `4·land − 2·adjacencies` |
| 3 | DFS Flood Fill | O(rows · cols) | O(rows · cols) | When the single island must be explored as a component |

---

## Approach 1 — Cell Scan, Count Exposed Sides

### Intuition

Every land cell is a unit square with 4 sides. A side is part of the island's perimeter precisely when the cell on the other side of it is **not** land — that is, water, or off the edge of the grid. So for each land cell, examine its 4 orthogonal neighbours and add 1 for every neighbour that is water or out of bounds.

### Algorithm

1. Initialise `perimeter = 0`.
2. For each cell `(r, c)` with `grid[r][c] == 1`:
   - For each of the 4 directions, if the neighbour is out of bounds or equals `0`, increment `perimeter`.
3. Return `perimeter`.

### Complexity

- **Time:** O(rows · cols) — a fixed 4 neighbour checks per cell.
- **Space:** O(1) — only a counter.

### Code

```go
func cellScan(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	// The four orthogonal neighbour offsets: up, down, left, right.
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	perimeter := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != 1 {
				continue // water contributes nothing
			}
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1] // neighbour coordinates
				// A side is on the perimeter if the neighbour is off-grid or water.
				if nr < 0 || nr >= rows || nc < 0 || nc >= cols || grid[nr][nc] == 0 {
					perimeter++
				}
			}
		}
	}
	return perimeter
}
```

### Dry Run

Example 3 is trivial; trace **Example 1** partially. `grid`:

```
row 0:  0 1 0 0
row 1:  1 1 1 0
row 2:  0 1 0 0
row 3:  1 1 0 0
```

Per land cell, count neighbours that are water/off-grid (its exposed sides):

| Cell (r,c) | up | down | left | right | exposed sides |
|------------|----|----|----|----|---------------|
| (0,1) | edge | land | water | water | 3 |
| (1,0) | water | water | edge | land | 3 |
| (1,1) | land | land | land | land | 0 |
| (1,2) | water | water | land | water | 3 |
| (2,1) | land | water | water | water | 3 |
| (3,0) | water | edge | edge | land | 3 |
| (3,1) | land | edge | land | water | 2 (down=edge, right=water) |

Wait — recount (3,1): up=(2,1) land, down=edge (+1), left=(3,0) land, right=(3,2) water (+1) → 2.

Sum: `3 + 3 + 0 + 3 + 3 + 3 + 2 = 17`? Recheck (3,0): up=(2,0) water (+1), down=edge (+1), left=edge (+1), right=(3,1) land → 3. And (0,1): up=edge(+1), down=(1,1) land, left=(0,0) water(+1), right=(0,2) water(+1) → 3.

Total = 3+3+0+3+3+3+2 = **17** — but expected is 16, so one row is miscounted above. Correct tally per cell (verified by the program): (0,1)=3, (1,0)=3, (1,1)=0, (1,2)=3, (2,1)=2, (3,0)=3, (3,1)=2 → sum **16**. (Cell (2,1): up=(1,1) land, down=(3,1) land, left=(2,0) water +1, right=(2,2) water +1 → 2.)

| Cell (r,c) | exposed sides |
|------------|---------------|
| (0,1) | 3 |
| (1,0) | 3 |
| (1,1) | 0 |
| (1,2) | 3 |
| (2,1) | 2 |
| (3,0) | 3 |
| (3,1) | 2 |

Sum = `3 + 3 + 0 + 3 + 2 + 3 + 2 = 16`. Result: `16` ✔

---

## Approach 2 — Contribution Counting (Optimal)

### Intuition

Give every land cell its full 4 sides, then correct for shared edges. Each pair of adjacent land cells shares exactly one internal edge; that edge was counted twice (once per cell) yet belongs to neither's perimeter, so subtract 2 for it. To count each shared edge exactly once, only look **up** and **left** from each cell (the down/right neighbours count the same edge from their side). Result: `perimeter = 4·land − 2·shared`, in a single pass with no direction array.

### Algorithm

1. `land = 0`, `shared = 0`.
2. For each cell `(r, c)` with `grid[r][c] == 1`:
   - `land++`.
   - If `r > 0` and the cell above is land, `shared++`.
   - If `c > 0` and the cell to the left is land, `shared++`.
3. Return `4·land − 2·shared`.

### Complexity

- **Time:** O(rows · cols) — one pass, constant work per cell.
- **Space:** O(1).

### Code

```go
func contributionCount(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	land, shared := 0, 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != 1 {
				continue
			}
			land++ // this cell brings 4 sides to the table
			if r > 0 && grid[r-1][c] == 1 {
				shared++ // shares an edge with the land cell above
			}
			if c > 0 && grid[r][c-1] == 1 {
				shared++ // shares an edge with the land cell to the left
			}
		}
	}
	// Each shared edge was double-counted, and it is not part of the perimeter.
	return 4*land - 2*shared
}
```

### Dry Run

Example 1: there are **7** land cells. Count adjacencies by looking up/left only:

| Cell (r,c) | up is land? | left is land? | shared added |
|------------|-------------|---------------|--------------|
| (0,1) | — (r=0) | (0,0)=water | 0 |
| (1,0) | (0,0)=water | — (c=0) | 0 |
| (1,1) | (0,1)=land | (1,0)=land | 2 |
| (1,2) | (0,2)=water | (1,1)=land | 1 |
| (2,1) | (1,1)=land | (2,0)=water | 1 |
| (3,0) | (2,0)=water | — (c=0) | 0 |
| (3,1) | (2,1)=land | (3,0)=land | 2 |

`land = 7`, `shared = 0+0+2+1+1+0+2 = 6`.

`perimeter = 4·7 − 2·6 = 28 − 12 = 16`. Result: `16` ✔

---

## Approach 3 — DFS Flood Fill (Count Boundary Edges)

### Intuition

Since there is exactly one island, a DFS from any land cell walks the entire island. As the DFS crosses a cell, count its boundary-contributing sides: stepping off the grid or into water means the side we just crossed is a boundary edge (`+1`); stepping into unvisited land recurses; stepping into already-visited land adds nothing (that internal edge is handled once). A `visited` grid prevents infinite recursion. This is the graph-search framing — the same template used to *count* or *measure* islands.

### Algorithm

1. Find any land cell and DFS from it.
2. In `dfs(r, c)`:
   - If `(r, c)` is off-grid or water → return `1` (a boundary edge).
   - If already visited → return `0`.
   - Otherwise mark visited and return the sum of `dfs` over the 4 neighbours.
3. The DFS return value is the perimeter.

### Complexity

- **Time:** O(rows · cols) — each cell entered once; neighbour lookups are O(1).
- **Space:** O(rows · cols) — the `visited` grid plus the recursion stack.

### Code

```go
func dfsPerimeter(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	visited := make([][]bool, rows) // guards against revisiting land cells
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	var dfs func(r, c int) int
	dfs = func(r, c int) int {
		// Off-grid or water: the side we crossed to get here is a boundary edge.
		if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] == 0 {
			return 1
		}
		if visited[r][c] {
			return 0 // internal shared edge, already accounted for
		}
		visited[r][c] = true // mark before recursing to avoid cycles
		// Sum boundary edges contributed through each of the 4 sides.
		return dfs(r-1, c) + dfs(r+1, c) + dfs(r, c-1) + dfs(r, c+1)
	}

	// Launch DFS from the first land cell (exactly one island guaranteed).
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 1 {
				return dfs(r, c)
			}
		}
	}
	return 0 // no land at all (not possible under constraints)
}
```

### Dry Run

Example 2: `grid = [[1]]` (a single land cell). DFS launches at `(0,0)`.

| Call | (r,c) | state | returns |
|------|-------|-------|---------|
| dfs(0,0) | (0,0) | land, unvisited → mark, recurse 4 sides | sum below |
| ↳ dfs(-1,0) | up | off-grid | 1 |
| ↳ dfs(1,0) | down | off-grid | 1 |
| ↳ dfs(0,-1) | left | off-grid | 1 |
| ↳ dfs(0,1) | right | off-grid | 1 |

`dfs(0,0) = 1 + 1 + 1 + 1 = 4`. Result: `4` ✔

---

## Key Takeaways

- **Perimeter = local edge accounting.** Two ways to see it: (a) sum each land cell's water/edge-facing sides, or (b) `4·land − 2·(adjacent land pairs)`. Both are O(cells).
- **Count each shared edge once by looking only up/left.** A recurring trick to avoid double-counting undirected adjacencies in a grid scan.
- **"Exactly one island" ⇒ a single DFS suffices.** No outer component loop needed; return as soon as the first land cell launches the search.
- The same grid + 4-neighbour skeleton powers #200 (Number of Islands), #695 (Max Area of Island), and #733 (Flood Fill) — only the per-cell bookkeeping changes.

---

## Related Problems

- LeetCode #200 — Number of Islands (count connected components)
- LeetCode #695 — Max Area of Island (size instead of perimeter)
- LeetCode #733 — Flood Fill (same traversal skeleton)
- LeetCode #1020 — Number of Enclaves (land not touching the border)
- LeetCode #733 — Flood Fill
