# 0417 — Pacific Atlantic Water Flow

> LeetCode #417 · Difficulty: Medium
> **Categories:** Array, Matrix, Depth-First Search, Breadth-First Search

---

## Problem Statement

There is an `m x n` rectangular island that borders both the **Pacific Ocean** and **Atlantic Ocean**. The **Pacific Ocean** touches the island's left and top edges, and the **Atlantic Ocean** touches the island's right and bottom edges.

The island is partitioned into a grid of square cells. You are given an `m x n` integer matrix `heights` where `heights[r][c]` represents the **height above sea level** of the cell at coordinate `(r, c)`.

The island receives a lot of rain, and the rain water can flow to neighboring cells directly north, south, east, and west if the neighboring cell's height is **less than or equal to** the current cell's height. Water can flow from any cell adjacent to an ocean into the ocean.

Return *a **2D list** of grid coordinates* `result` *where* `result[i] = [ri, ci]` *denotes that rain water can flow from cell* `(ri, ci)` *to **both** the Pacific and Atlantic oceans.*

**Example 1:**

```
Input: heights = [[1,2,2,3,5],[3,2,3,4,4],[2,4,5,3,1],[6,7,1,4,5],[5,1,1,2,4]]
Output: [[0,4],[1,3],[1,4],[2,2],[3,0],[3,1],[4,0]]
Explanation: The following cells can flow to the Pacific and Atlantic oceans, as shown below:
[0,4]: [0,4] -> Pacific Ocean
       [0,4] -> Atlantic Ocean
[1,3]: [1,3] -> [0,3] -> Pacific Ocean
       [1,3] -> [1,4] -> Atlantic Ocean
[1,4]: [1,4] -> [1,3] -> [0,3] -> Pacific Ocean
       [1,4] -> Atlantic Ocean
[2,2]: [2,2] -> [1,2] -> [0,2] -> Pacific Ocean
       [2,2] -> [2,3] -> [2,4] -> Atlantic Ocean
[3,0]: [3,0] -> Pacific Ocean
       [3,0] -> [4,0] -> Atlantic Ocean
[3,1]: [3,1] -> [3,0] -> Pacific Ocean
       [3,1] -> [4,1] -> Atlantic Ocean
[4,0]: [4,0] -> Pacific Ocean
       [4,0] -> Atlantic Ocean
Note that there are other possible paths for these cells to flow to the Pacific and Atlantic oceans.
```

**Example 2:**

```
Input: heights = [[1]]
Output: [[0,0]]
Explanation: The water can flow from the only cell to the Pacific and Atlantic oceans.
```

**Constraints:**

- `m == heights.length`
- `n == heights[r].length`
- `1 <= m, n <= 200`
- `0 <= heights[r][c] <= 10^5`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Graph BFS/DFS on a grid** — cells are nodes, "flows to" edges connect a cell to lower/equal neighbours; the core trick is flooding *inward from the ocean borders* by reversing the edge direction → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Matrix traversal** — 4-directional movement with bounds checking over a 2D grid, plus two overlaid reachability masks whose intersection is the answer → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (flood from every cell) | O((m·n)²) | O(m·n) | Conceptual baseline; 200×200 grid → up to 1.6·10⁹ work, risky |
| 2 | Reverse DFS from ocean borders (Optimal) | O(m·n) | O(m·n) | The canonical answer; concise recursion |
| 3 | Reverse BFS from ocean borders | O(m·n) | O(m·n) | Same complexity, iterative queue — safe against deep recursion |

---

## Approach 1 — Brute Force (BFS/DFS From Every Cell)

### Intuition

Take the problem statement literally. A cell drains to an ocean iff there is a path of non-increasing heights from it to that ocean's border. So for **each** cell, launch a flood that only steps to equal-or-lower neighbours, and record whether it ever touches a Pacific edge (top/left) and whether it touches an Atlantic edge (bottom/right). The cell is in the answer when both flags fire. Correct but wasteful: every cell pays for its own full-grid search.

### Algorithm

1. For each cell `(r, c)`:
   - DFS outward, moving to any in-bounds neighbour with `height <= current height`.
   - When the flood lands on `r == 0 || c == 0`, set `pacific = true`; when it lands on `r == m-1 || c == n-1`, set `atlantic = true`.
2. If both flags are set, append `[r, c]` to the result.

### Complexity

- **Time:** O((m·n)²) — an O(m·n) flood repeated for each of the m·n cells.
- **Space:** O(m·n) — a fresh visited grid and recursion stack per flood.

### Code

```go
func bruteForce(heights [][]int) [][]int {
	m, n := len(heights), len(heights[0])
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // N, S, W, E
	var result [][]int

	// canReachBoth floods from (sr, sc) and returns whether the flood touches
	// each ocean's border.
	canReachBoth := func(sr, sc int) (pacific, atlantic bool) {
		visited := make([][]bool, m)
		for i := range visited {
			visited[i] = make([]bool, n)
		}
		var dfs func(r, c int)
		dfs = func(r, c int) {
			visited[r][c] = true
			if r == 0 || c == 0 { // top or left edge → Pacific
				pacific = true
			}
			if r == m-1 || c == n-1 { // bottom or right edge → Atlantic
				atlantic = true
			}
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr < 0 || nr >= m || nc < 0 || nc >= n {
					continue // off the board
				}
				if visited[nr][nc] {
					continue // already flooded
				}
				if heights[nr][nc] > heights[r][c] {
					continue // uphill — water can't flow there
				}
				dfs(nr, nc) // step downhill/flat
			}
		}
		dfs(sr, sc)
		return
	}

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if p, a := canReachBoth(r, c); p && a {
				result = append(result, []int{r, c})
			}
		}
	}
	return result
}
```

### Dry Run

Example 1, focusing on a few cells (grid is 5×5). Flood moves only to equal/lower heights.

| Cell | height | Reaches Pacific edge? | Reaches Atlantic edge? | In answer? |
|------|--------|-----------------------|------------------------|------------|
| (0,4)=5 | 5 | yes — it *is* on the top edge | yes — it *is* on the right edge | ✔ |
| (2,2)=5 | 5 | 5→(1,2)3→(0,2)2 top edge | 5→(2,3)3→(2,4)1 right edge | ✔ |
| (0,0)=1 | 1 | on top+left edge (Pacific) | 1→ can only reach ≤1 cells; can't get downhill to bottom/right | ✘ |
| (3,3)=4 | 4 | cannot climb up to top/left | 4→(2,3)3? no that's uphill-blocked… drains right/bottom | ✘ (not in output) |

Collecting all cells that satisfy both → `[[0,4],[1,3],[1,4],[2,2],[3,0],[3,1],[4,0]]` ✔

---

## Approach 2 — Reverse DFS From Ocean Borders (Optimal)

### Intuition

Flip the question. Rather than "can water leave this cell for the ocean?", ask "which cells can send water **to** this ocean?" and start from the ocean itself. Forward flow goes to equal-or-lower neighbours; **reversed** flow therefore goes to equal-or-**higher** neighbours (a cell can feed us only if it sits at our height or above). Run one climbing flood seeded from every Pacific-border cell (top row + left column) to mark all cells draining to the Pacific, and another seeded from every Atlantic-border cell (bottom row + right column). Cells marked by **both** floods are the answer. Two grid-sized passes replace the brute force's m·n passes.

### Algorithm

1. Build two boolean grids `pacific` and `atlantic`.
2. `dfs(r, c, ocean)`: mark `ocean[r][c]`; recurse into every neighbour with `height >= heights[r][c]` that isn't marked yet.
3. Seed `pacific` from the whole top row and left column; seed `atlantic` from the whole bottom row and right column.
4. Every `(r, c)` with `pacific[r][c] && atlantic[r][c]` goes into the result.

### Complexity

- **Time:** O(m·n) — each cell is marked at most once per ocean, so ≤ 2·m·n visits.
- **Space:** O(m·n) — two reachability grids plus recursion depth.

### Code

```go
func reverseDFS(heights [][]int) [][]int {
	m, n := len(heights), len(heights[0])
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	pacific := make([][]bool, m)  // cells that can drain to the Pacific
	atlantic := make([][]bool, m) // cells that can drain to the Atlantic
	for i := 0; i < m; i++ {
		pacific[i] = make([]bool, n)
		atlantic[i] = make([]bool, n)
	}

	// dfs climbs inward: from (r,c) visit neighbours with height >= heights[r][c].
	var dfs func(r, c int, ocean [][]bool)
	dfs = func(r, c int, ocean [][]bool) {
		ocean[r][c] = true // this cell drains to the ocean we started from
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr < 0 || nr >= m || nc < 0 || nc >= n {
				continue // off the board
			}
			if ocean[nr][nc] {
				continue // already reached
			}
			if heights[nr][nc] < heights[r][c] {
				continue // reverse flow must go UPHILL-or-flat; lower neighbour can't feed us
			}
			dfs(nr, nc, ocean) // climb to the equal/higher neighbour
		}
	}

	// Seed the Pacific from the top row and left column.
	for c := 0; c < n; c++ {
		dfs(0, c, pacific) // top edge
	}
	for r := 0; r < m; r++ {
		dfs(r, 0, pacific) // left edge
	}
	// Seed the Atlantic from the bottom row and right column.
	for c := 0; c < n; c++ {
		dfs(m-1, c, atlantic) // bottom edge
	}
	for r := 0; r < m; r++ {
		dfs(r, n-1, atlantic) // right edge
	}

	var result [][]int
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if pacific[r][c] && atlantic[r][c] { // intersection of both floods
				result = append(result, []int{r, c})
			}
		}
	}
	return result
}
```

### Dry Run

Example 1. Seed the Pacific flood from the top row `[1,2,2,3,5]` and left column `[1,3,2,6,5]`, climbing to `>=` neighbours.

| Phase | Seed cell | Climbs to (height ≥ current) | Marks (sample) |
|-------|-----------|------------------------------|----------------|
| Pacific | (0,0)=1 | (0,1)=2, (1,0)=3, … | top row, most of rows 0–1 |
| Pacific | (3,0)=6 | (4,0)=5? no (5<6 blocked from 6) but (3,0) itself marked; (4,0) seeded via left col | (3,0),(4,0) |
| Atlantic | (4,4)=4 | (4,3)=2? no; (3,4)=5 ≥4 yes | right/bottom region |
| Atlantic | (2,4)=1 | (1,4)=4 ≥1, (2,3)=3 ≥1 … | (1,3),(1,4),(2,2) reachable |

Intersection of the two masks: `{(0,4),(1,3),(1,4),(2,2),(3,0),(3,1),(4,0)}`. Result: `[[0,4],[1,3],[1,4],[2,2],[3,0],[3,1],[4,0]]` ✔

---

## Approach 3 — Reverse BFS From Ocean Borders

### Intuition

The reverse-flood idea again, but with an explicit FIFO queue instead of recursion. Seed one queue with every border cell of an ocean (all marked), then repeatedly pop a cell and enqueue any equal-or-higher unvisited neighbour — those are exactly the cells that could send water down into the current one. BFS and DFS mark the identical set; the queue simply guarantees bounded stack usage, which matters when a 200×200 grid could otherwise recurse 40 000 deep.

### Algorithm

1. Collect border seeds: Pacific = top row + left column; Atlantic = bottom row + right column.
2. `bfs(starts)`: mark all seeds, enqueue them; pop cells and for each neighbour with `height >= current`, mark and enqueue.
3. Intersect the two returned grids into the result.

### Complexity

- **Time:** O(m·n) — every cell enters each ocean's queue at most once.
- **Space:** O(m·n) — reachability grids plus the queue (≤ m·n entries).

### Code

```go
func reverseBFS(heights [][]int) [][]int {
	m, n := len(heights), len(heights[0])
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	// bfs marks every cell that can reach the ocean whose border cells seed it.
	bfs := func(starts [][2]int) [][]bool {
		reach := make([][]bool, m)
		for i := range reach {
			reach[i] = make([]bool, n)
		}
		queue := make([][2]int, 0, len(starts))
		for _, s := range starts {
			reach[s[0]][s[1]] = true // border cells trivially drain to their ocean
			queue = append(queue, s)
		}
		for len(queue) > 0 {
			cell := queue[0] // dequeue front (FIFO)
			queue = queue[1:]
			r, c := cell[0], cell[1]
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr < 0 || nr >= m || nc < 0 || nc >= n {
					continue
				}
				if reach[nr][nc] {
					continue
				}
				if heights[nr][nc] < heights[r][c] {
					continue // neighbour lower → cannot flow uphill into us
				}
				reach[nr][nc] = true // neighbour drains to this ocean too
				queue = append(queue, [2]int{nr, nc})
			}
		}
		return reach
	}

	// Border seeds for each ocean.
	var pacStarts, atlStarts [][2]int
	for c := 0; c < n; c++ {
		pacStarts = append(pacStarts, [2]int{0, c})     // top row
		atlStarts = append(atlStarts, [2]int{m - 1, c}) // bottom row
	}
	for r := 0; r < m; r++ {
		pacStarts = append(pacStarts, [2]int{r, 0})     // left col
		atlStarts = append(atlStarts, [2]int{r, n - 1}) // right col
	}

	pacific := bfs(pacStarts)
	atlantic := bfs(atlStarts)

	var result [][]int
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if pacific[r][c] && atlantic[r][c] {
				result = append(result, []int{r, c})
			}
		}
	}
	return result
}
```

### Dry Run

Example 1, Pacific queue. Seeds = top row + left column, all marked.

| Step | Popped | height | Enqueues (neighbour height ≥ popped) |
|------|--------|--------|--------------------------------------|
| 1 | (0,0) | 1 | (already seeded) |
| 2 | (0,1) | 2 | (1,1)=2 ≥2 → mark, enqueue |
| 3 | (1,0) | 3 | (1,1) already; (2,0)=2 <3 skip |
| … | (1,1) | 2 | (1,2)=3 ≥2 → enqueue; (2,1)=4 ≥2 → enqueue |
| … | (2,1) | 4 | (2,2)=5 ≥4 → enqueue |

The Pacific mask grows to include `(2,2)`. The Atlantic BFS (seeded from bottom row + right column) also reaches `(2,2)` via `1→3→5` uphill from `(2,4)`. Intersection over the whole grid → `[[0,4],[1,3],[1,4],[2,2],[3,0],[3,1],[4,0]]` ✔

---

## Key Takeaways

- **Reverse the flow.** "Which sources reach a sink?" over many sinks is often far cheaper solved as "which nodes does each sink reach?" seeded from the sinks. Here it turns O((mn)²) into O(mn). Forward edge `A→B` when `h[B] ≤ h[A]` becomes reverse edge `B→A` when `h[A] ≥ h[B]` — i.e. climb non-decreasing heights.
- **Multi-source flood.** Seeding a single BFS/DFS with *all* border cells at once computes an entire reachability region in one pass — no need to loop per source.
- **Intersection of two masks.** Independent reachability sets, ANDed cell-by-cell, is a clean pattern for "reaches both X and Y" grid problems.
- **DFS vs BFS are interchangeable here** — same visited-set semantics; pick BFS (explicit queue) when recursion depth (up to m·n ≈ 40 000) is a stack-overflow risk.

---

## Related Problems

- LeetCode #200 — Number of Islands (grid flood fill)
- LeetCode #130 — Surrounded Regions (flood inward from the border — same reversal trick)
- LeetCode #542 — 01 Matrix (multi-source BFS from all zeros)
- LeetCode #994 — Rotting Oranges (multi-source BFS over a grid)
- LeetCode #1020 — Number of Enclaves (border-seeded flood fill)
