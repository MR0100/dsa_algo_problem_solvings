# 0317 — Shortest Distance from All Buildings

> LeetCode #317 · Difficulty: Hard
> **Categories:** Array, Matrix, Breadth-First Search, Multi-Source BFS

---

## Problem Statement

You are given an `m x n` grid `grid` of values `0`, `1`, or `2`, where:

- each `0` marks an **empty land** that you can pass by freely,
- each `1` marks a **building** that you cannot pass through, and
- each `2` marks an **obstacle** that you cannot pass through.

You want to build a house on an empty land that reaches all buildings in the
**shortest total travel distance**. You can only move up, down, left, and right.

Return the shortest travel distance for such a house. If it is not possible to
build such a house according to the above rules, return `-1`.

The **total travel distance** is the sum of the distances between the house and
each of the buildings, using the Manhattan-style shortest walkable path.

**Example 1:**

```
Input:  grid = [[1,0,2,0,1],[0,0,0,0,0],[0,0,1,0,0]]
Output: 7
```

Explanation: Given three buildings at `(0,0)`, `(0,4)`, `(2,2)`, and an obstacle
at `(0,2)`. The point `(1,2)` is an ideal empty land to build a house, as the
total travel distance of `3 + 3 + 1 = 7` is minimal. So return `7`.

**Example 2:**

```
Input:  grid = [[1,0]]
Output: 1
```

**Example 3:**

```
Input:  grid = [[1]]
Output: -1
```

**Constraints:**

- `m == grid.length`
- `n == grid[i].length`
- `1 <= m, n <= 50`
- `grid[i][j]` is either `0`, `1`, or `2`.
- There will be **at least one** building in the grid.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★★ Very High  | 2024          |
| Amazon    | ★★★★☆ High       | 2023          |
| Meta      | ★★★☆☆ Medium     | 2023          |
| Uber      | ★★★☆☆ Medium     | 2022          |
| Microsoft | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Breadth-First Search (multi-source style)** — BFS gives shortest step counts
  on an unweighted grid; we run it once per building →
  see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Matrix Traversal** — 4-directional movement with bounds/obstacle checks →
  see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Queue / Deque** — the BFS frontier is a FIFO queue →
  see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS from every empty cell | O((m·n)²) | O(m·n) | Small grids; most intuitive |
| 2 | BFS from every building (Optimal) | O(B·m·n) | O(m·n) | Buildings ≪ empties; the intended solution |

---

## Approach 1 — BFS From Every Empty Cell

### Intuition
The house must sit on an empty cell. For a fixed empty cell, one BFS outward
computes the shortest distance to every building simultaneously. Sum those
distances; if the BFS cannot reach some building, this cell is invalid. Try all
empty cells and keep the minimum.

### Algorithm
1. Count total buildings `B`.
2. For each empty cell `(i,j)`:
   - BFS outward, only walking through empty cells (`0`); obstacles (`2`) and
     buildings (`1`) are not walked through.
   - When a building is first reached at layer `dist`, add `dist` to a running
     `totalDist` and increment `reached`.
3. If `reached == B`, update the best answer with `totalDist`.
4. Return the best, or `-1` if no cell reaches all buildings.

### Complexity
- **Time:** O((m·n)²) — up to `m·n` empty starts, each an O(m·n) BFS.
- **Space:** O(m·n) — a fresh `visited` matrix and the BFS queue per start.

### Code
```go
func bfsFromEmpty(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return -1
	}
	n := len(grid[0])
	totalBuildings := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				totalBuildings++
			}
		}
	}
	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	best := -1
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] != 0 {
				continue
			}
			visited := make([][]bool, m)
			for r := range visited {
				visited[r] = make([]bool, n)
			}
			queue := [][2]int{{i, j}}
			visited[i][j] = true
			dist := 0
			totalDist := 0
			reached := 0
			for len(queue) > 0 {
				dist++
				next := [][2]int{}
				for _, cell := range queue {
					for _, d := range dirs {
						nr, nc := cell[0]+d[0], cell[1]+d[1]
						if nr < 0 || nr >= m || nc < 0 || nc >= n || visited[nr][nc] {
							continue
						}
						visited[nr][nc] = true
						if grid[nr][nc] == 1 {
							totalDist += dist
							reached++
						} else if grid[nr][nc] == 0 {
							next = append(next, [2]int{nr, nc})
						}
					}
				}
				queue = next
			}
			if reached == totalBuildings {
				if best == -1 || totalDist < best {
					best = totalDist
				}
			}
		}
	}
	return best
}
```

### Dry Run
Input `grid2 = [[1,0]]`, B = 1.

| Step | Detail |
|------|--------|
| Empty cells | only `(0,1)` |
| BFS from `(0,1)` | layer dist=1 visits `(0,0)` which is a building |
| totalDist | `+= 1` → 1; reached = 1 |
| reached == B? | 1 == 1 → valid |
| best | updated to 1 |

Return `1`. ✓ (For `grid3 = [[1]]` there are no empty cells, so best stays `-1`.)

---

## Approach 2 — BFS From Every Building (Optimal)

### Intuition
Buildings are typically far fewer than empty cells, so BFS from buildings is
cheaper. For each building we add its shortest distance into a `total[r][c]`
accumulator for every reachable empty cell. To require an empty cell be reached
by **all** buildings without repeatedly rebuilding a visited matrix, we use a
sliding "marker": in round `k` (starting 0) a cell is walkable only if its grid
value equals `emptyMarker = -k`. After visiting we decrement it, so only cells
touched by every building so far remain walkable in the next round. At the end,
cells whose value equals the final `emptyMarker` are exactly those reachable by
all buildings; the answer is the minimum `total` among them.

### Algorithm
1. Copy the grid; `emptyMarker = 0`; `total[r][c] = 0`.
2. For each building:
   - BFS. A neighbour is walkable iff `g[nr][nc] == emptyMarker`.
   - On visiting: `g[nr][nc]--`, and `total[nr][nc] += dist`.
   - After the BFS: `emptyMarker--`.
3. Answer = minimum `total[r][c]` over cells with `g[r][c] == emptyMarker` and
   `total > 0`; `-1` if none.

### Complexity
- **Time:** O(B·m·n) — one O(m·n) BFS per building `B`.
- **Space:** O(m·n) — the `total` accumulator, grid copy, and BFS queue.

### Code
```go
func bfsFromBuildings(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return -1
	}
	n := len(grid[0])
	total := make([][]int, m)
	for i := range total {
		total[i] = make([]int, n)
	}
	g := make([][]int, m)
	for i := range g {
		g[i] = make([]int, n)
		copy(g[i], grid[i])
	}
	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	emptyMarker := 0
	best := -1
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if g[i][j] != 1 {
				continue
			}
			queue := [][2]int{{i, j}}
			dist := 0
			for len(queue) > 0 {
				dist++
				next := [][2]int{}
				for _, cell := range queue {
					for _, d := range dirs {
						nr, nc := cell[0]+d[0], cell[1]+d[1]
						if nr < 0 || nr >= m || nc < 0 || nc >= n || g[nr][nc] != emptyMarker {
							continue
						}
						g[nr][nc]--
						total[nr][nc] += dist
						next = append(next, [2]int{nr, nc})
					}
				}
				queue = next
			}
			emptyMarker--
		}
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if g[i][j] == emptyMarker && total[i][j] > 0 {
				if best == -1 || total[i][j] < best {
					best = total[i][j]
				}
			}
		}
	}
	return best
}
```

### Dry Run
Input `grid2 = [[1,0]]`. One building at `(0,0)`.

| Phase | State |
|-------|-------|
| Init | `emptyMarker = 0`, `g = [[1,0]]`, `total = [[0,0]]` |
| BFS from `(0,0)` | dist=1: neighbour `(0,1)` has `g=0 == emptyMarker` → walkable |
| Visit `(0,1)` | `g[0][1] = -1`, `total[0][1] += 1` → 1 |
| After round | `emptyMarker = -1` |
| Final scan | `(0,1)`: `g == -1 == emptyMarker` and total=1>0 → best = 1 |

Return `1`. ✓ (For `grid3 = [[1]]` the BFS visits nothing; no cell has
`total>0`, so best stays `-1`.)

---

## Key Takeaways

- **Choose the smaller source set.** BFS-from-buildings is O(B·m·n); when
  buildings are sparse this crushes the O((m·n)²) BFS-from-empties.
- **Sliding marker trick.** Decrementing the grid value each round both records
  "reached by all previous buildings" and doubles as a visited flag — no
  per-round visited matrix reset needed.
- **Unweighted grid ⇒ BFS layers = shortest distance.** Increment `dist` once
  per BFS ring, not per cell.
- **Validity is a reachability constraint**: a candidate is only valid if every
  building can reach it, hence the "reached by all" bookkeeping.

---

## Related Problems

- LeetCode #286 — Walls and Gates (multi-source BFS)
- LeetCode #542 — 01 Matrix (multi-source BFS distance)
- LeetCode #994 — Rotting Oranges (multi-source BFS)
- LeetCode #675 — Cut Off Trees for Golf Event (BFS shortest path on grid)
