# 0200 — Number of Islands

> LeetCode #200 · Difficulty: Medium
> **Categories:** Array, Depth-First Search, Breadth-First Search, Union Find, Matrix

---

## Problem Statement

Given an `m x n` 2D binary grid `grid` which represents a map of `'1'`s (land) and `'0'`s (water), return *the number of islands*.

An **island** is surrounded by water and is formed by connecting adjacent lands horizontally or vertically. You may assume all four edges of the grid are all surrounded by water.

**Example 1:**

```
Input: grid = [
  ["1","1","1","1","0"],
  ["1","1","0","1","0"],
  ["1","1","0","0","0"],
  ["0","0","0","0","0"]
]
Output: 1
```

**Example 2:**

```
Input: grid = [
  ["1","1","0","0","0"],
  ["1","1","0","0","0"],
  ["0","0","1","0","0"],
  ["0","0","0","1","1"]
]
Output: 3
```

**Constraints:**

- `m == grid.length`
- `n == grid[i].length`
- `1 <= m, n <= 300`
- `grid[i][j]` is `'0'` or `'1'`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Depth-First Search / Breadth-First Search on a grid** — the grid is an implicit graph whose vertices are cells and whose edges join orthogonal neighbours; counting islands is counting connected components, and a flood fill from each unvisited land cell erases one whole component → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Matrix Traversal** — the outer scan sweeps every `(r, c)` and the flood fill moves by the four orthogonal offsets `{(±1,0),(0,±1)}` with bounds checks; the standard grid-neighbour pattern → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Queue / Deque** — BFS uses a FIFO frontier queue, expanding an island ring by ring instead of recursively → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Union-Find (Disjoint Set Union)** — every land cell starts as its own component and adjacent land pairs are merged; the surviving component count is the answer, and this is the version that generalises to the dynamic follow-up (#305) → see [`/dsa/union_find.md`](/dsa/union_find.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DFS Flood Fill | O(m·n) | O(m·n) | Fewest lines; the default answer. Watch recursion depth on a 300×300 all-land grid |
| 2 | BFS Flood Fill | O(m·n) | O(min(m,n)) | Same idea, iterative — safe when a deep recursion stack is a concern |
| 3 | Union-Find (DSU) | O(m·n · α(m·n)) | O(m·n) | When land is added incrementally (follow-up #305) or you want a components framework |

---

## Approach 1 — DFS Flood Fill

### Intuition

Scan every cell in reading order. The first time the scan lands on a `'1'`, that cell belongs to an island nobody has counted yet — so bump the counter by one, then **sink** the entire island reachable from it by flipping every connected `'1'` to `'0'`. Sinking doubles as the visited-mark: because the whole island is now water, the outer scan can never re-enter it. Therefore *islands counted = flood fills started*.

### Algorithm

1. Initialise `count = 0`.
2. For every cell `(r, c)` in row-major order: if `grid[r][c] == '1'`, increment `count` and call `sinkDFS(r, c)`.
3. `sinkDFS` returns immediately for out-of-bounds cells and for cells that are not `'1'`; otherwise it sets the cell to `'0'` and recurses into its four orthogonal neighbours (down, up, right, left).
4. Return `count`.

### Complexity

- **Time:** O(m·n) — the outer loop visits every cell once, and each cell is sunk by `sinkDFS` at most once, so total work is linear in the grid size.
- **Space:** O(m·n) — the recursion stack. In the worst case (the entire grid is land arranged as one snake), the DFS recurses through all `m·n` cells before unwinding.

### Code

```go
func dfsFloodFill(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}
	count := 0
	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] == '1' {
				count++             // a new, never-seen island starts here
				sinkDFS(grid, r, c) // erase it so it is counted exactly once
			}
		}
	}
	return count
}

// sinkDFS turns the whole island containing (r, c) into water.
func sinkDFS(grid [][]byte, r, c int) {
	// Stop at the grid border and at water / already-sunk cells.
	if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) || grid[r][c] != '1' {
		return
	}
	grid[r][c] = '0'      // sink: acts as the "visited" mark
	sinkDFS(grid, r+1, c) // down
	sinkDFS(grid, r-1, c) // up
	sinkDFS(grid, r, c+1) // right
	sinkDFS(grid, r, c-1) // left
}
```

### Dry Run

Example 1 grid (rows top-to-bottom, `.` shown for `'0'`):

```
row0: 1 1 1 1 .
row1: 1 1 . 1 .
row2: 1 1 . . .
row3: . . . . .
```

The outer scan reaches `(0,0)` first, which is `'1'`.

| Event | Cell | `count` | Grid state after event |
|-------|------|---------|-------------------------|
| Scan hits land at `(0,0)` | `(0,0)` | 1 | flood fill begins |
| `sinkDFS` sinks the connected component | all `'1'`s reachable from `(0,0)`: `(0,0)(0,1)(0,2)(0,3)(1,0)(1,1)(1,3)(2,0)(2,1)` | 1 | every `'1'` is now `'0'` → whole grid is water |
| Scan continues `(0,1)…(3,4)` | — | 1 | no cell is `'1'` anymore, no new fill |

Loop ends. Result: **1** ✔ — the four `'1'`s at `(0,3)` and `(1,3)` connect back to the corner block through `(0,2)→(0,3)`, so it is a single island.

---

## Approach 2 — BFS Flood Fill

### Intuition

The counting argument is identical to DFS — each unvisited land cell seeds exactly one island — but the island is explored with an **iterative FIFO queue** instead of the call stack. BFS expands the island ring by ring: enqueue the seed, then repeatedly pop a cell and enqueue its land neighbours. The key discipline is **mark-on-enqueue**: sink a neighbour to `'0'` the instant it is pushed, so no cell ever enters the queue twice. Prefer BFS when recursion depth is a worry — a 300×300 all-land grid would drive the DFS stack to 90,000 frames.

### Algorithm

1. Initialise `count = 0` and the four direction offsets.
2. Scan every cell `(r, c)`. When `grid[r][c] == '1'`: increment `count`, set it to `'0'`, and enqueue `(r, c)`.
3. While the queue is non-empty: dequeue a cell; for each in-bounds neighbour still equal to `'1'`, set it to `'0'` **and** enqueue it.
4. Return `count`.

### Complexity

- **Time:** O(m·n) — every cell is enqueued at most once (guarded by mark-on-enqueue), and each dequeue inspects four neighbours in constant time.
- **Space:** O(min(m,n)) — the frontier at any moment is roughly a diagonal band across the island; O(m·n) is the loose upper bound.

### Code

```go
func bfsFloodFill(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}
	rows, cols := len(grid), len(grid[0])
	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} // the four orthogonal moves
	count := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != '1' {
				continue // water or an already-sunk island cell
			}
			count++          // new island discovered
			grid[r][c] = '0' // mark before enqueueing to avoid duplicates
			queue := [][2]int{{r, c}}
			for len(queue) > 0 {
				cell := queue[0]
				queue = queue[1:]
				for _, d := range dirs {
					nr, nc := cell[0]+d[0], cell[1]+d[1]
					// Sink in-bounds land neighbours the moment we see them.
					if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == '1' {
						grid[nr][nc] = '0'
						queue = append(queue, [2]int{nr, nc})
					}
				}
			}
		}
	}
	return count
}
```

### Dry Run

Example 2 grid:

```
row0: 1 1 . . .
row1: 1 1 . . .
row2: . . 1 . .
row3: . . . 1 1
```

Scanning row-major, the three seeds fire one BFS each.

| Seed found | `count` | BFS sinks (cells enqueued & marked `'0'`) | Notes |
|------------|---------|--------------------------------------------|-------|
| `(0,0)` | 1 | `(0,0)→(1,0),(0,1)→(1,1)` | 2×2 block in the top-left, one island |
| `(2,2)` | 2 | `(2,2)` (no land neighbour) | a lone cell, one island |
| `(3,3)` | 3 | `(3,3)→(3,4)` | horizontal pair in bottom-right, one island |

After each BFS the island's cells are water, so the scan skips over them. No other `'1'` remains. Result: **3** ✔

---

## Approach 3 — Union-Find (Disjoint Set Union)

### Intuition

Restate the problem: "number of islands" is precisely "number of connected components" in the graph whose vertices are land cells and whose edges join orthogonal land neighbours. A DSU maintains a live component count under edge insertions. So: give every land cell its own component, then walk the grid once and `union` each land cell with its right and down land neighbours (those two directions cover every adjacent pair exactly once). The surviving component count is the answer. This is also the approach that adapts cleanly to the dynamic follow-up **#305 (Number of Islands II)**, where land appears one cell at a time.

### Algorithm

1. Build the DSU: flatten `(r, c)` to id `r*cols + c`; every id is its own parent, and `count` starts equal to the number of land cells.
2. For every land cell `(r, c)`: if the cell below `(r+1, c)` is land, `union` them; if the cell to the right `(r, c+1)` is land, `union` them. Each successful merge decrements `count`.
3. `find` uses path halving and `union` uses union-by-rank to keep the trees near-flat.
4. Return `count`.

### Complexity

- **Time:** O(m·n · α(m·n)) — one pass over the grid performing O(m·n) union/find operations, each near-constant; α is the inverse Ackermann function (≤ 5 for any realistic grid), so effectively linear.
- **Space:** O(m·n) — the `parent` and `rank` arrays, one entry per cell.

### Code

```go
func unionFind(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}
	rows, cols := len(grid), len(grid[0])
	d := newDSU(grid)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != '1' {
				continue // only land cells create edges
			}
			if r+1 < rows && grid[r+1][c] == '1' {
				d.union(r*cols+c, (r+1)*cols+c) // edge to the cell below
			}
			if c+1 < cols && grid[r][c+1] == '1' {
				d.union(r*cols+c, r*cols+c+1) // edge to the cell on the right
			}
		}
	}
	return d.count
}
```

Supporting DSU (path halving + union by rank):

```go
type dsu struct {
	parent []int // parent[i] = parent of i; i is a root when parent[i] == i
	rank   []int // rank[i] = upper bound on the height of the tree rooted at i
	count  int   // number of live components (islands not yet merged)
}

func (d *dsu) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]] // path halving: skip a generation
		x = d.parent[x]
	}
	return x
}

func (d *dsu) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return // already part of the same island
	}
	if d.rank[ra] < d.rank[rb] {
		ra, rb = rb, ra // attach the shorter tree under the taller one
	}
	d.parent[rb] = ra
	if d.rank[ra] == d.rank[rb] {
		d.rank[ra]++ // equal heights → merged tree grows by one
	}
	d.count-- // two islands fused into one
}
```

### Dry Run

Example 2 grid, cells flattened as `id = r*5 + c`:

```
row0: id0=1  id1=1  id2=.  id3=.  id4=.
row1: id5=1  id6=1  id7=.  id8=.  id9=.
row2: id10=. id11=. id12=1 id13=. id14=.
row3: id15=. id16=. id17=. id18=1 id19=1
```

Initial land cells: `{0,1,5,6,12,18,19}` → `count = 7`.

| Cell `(r,c)` | id | Down land? → union | Right land? → union | `count` after |
|--------------|----|--------------------|--------------------|---------------|
| `(0,0)` | 0 | `(1,0)` id5 → union(0,5) | `(0,1)` id1 → union(0,1) | 7→6→5 |
| `(0,1)` | 1 | `(1,1)` id6 → union(1,6) | — (id2 water) | 5→4 |
| `(1,0)` | 5 | — (id10 water) | `(1,1)` id6 → union(5,6) already joined | 4 (no change) |
| `(1,1)` | 6 | — | — | 4 |
| `(2,2)` | 12 | — (id17 water) | — (id13 water) | 4 |
| `(3,3)` | 18 | — (row 4 absent) | `(3,4)` id19 → union(18,19) | 4→3 |
| `(3,4)` | 19 | — | — | 3 |

Three components survive: `{0,1,5,6}`, `{12}`, `{18,19}`. Result: **3** ✔

---

## Key Takeaways

- **"Number of islands" = number of connected components** in the grid-as-graph. Recognising this reframing unlocks DFS, BFS, and Union-Find as three interchangeable tools for the same shape of problem.
- **Sink as you visit.** Flipping `'1' → '0'` reuses the input grid as the visited set — no extra `visited` array, O(1) extra bookkeeping. If mutating the input is disallowed, deep-copy first (as `cloneGrid` does here) or keep a separate boolean matrix.
- **Mark-on-enqueue, not on-dequeue, in BFS.** Marking a cell the moment it is pushed prevents the same cell from being queued by two different neighbours — a classic bug that inflates memory and can double-count.
- **DFS vs BFS is a stack-depth trade.** Both are O(m·n) time; DFS is shorter to write but can blow the recursion stack on huge components, where BFS's explicit queue is safer.
- **Reach for Union-Find when connectivity is dynamic.** For a static grid, flood fill is simpler; but when land is added incrementally (#305), DSU answers each addition in near-constant amortised time, which flood fill cannot.

---

## Related Problems

- LeetCode #305 — Number of Islands II (dynamic land additions — the Union-Find follow-up)
- LeetCode #695 — Max Area of Island (same flood fill, track component size)
- LeetCode #463 — Island Perimeter (grid traversal, count exposed edges)
- LeetCode #130 — Surrounded Regions (flood fill from the border)
- LeetCode #547 — Number of Provinces (connected components via Union-Find / DFS)
- LeetCode #200 variants #733 — Flood Fill (the core sink primitive on its own)
- LeetCode #1254 — Number of Closed Islands (flood fill with border exclusion)
