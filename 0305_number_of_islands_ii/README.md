# 0305 — Number of Islands II

> LeetCode #305 · Difficulty: Hard (Premium / Locked)
> **Categories:** Array, Union Find, Graph

---

## Problem Statement

You are given an empty 2D binary grid `grid` of size `m x n`. The grid represents a map where `0`'s represent water and `1`'s represent land. Initially, all the cells of `grid` are water cells (i.e., all the cells are `0`'s).

We may perform an **add land** operation which turns the water at position into a land. You are given an array `positions` where `positions[i] = [ri, ci]` is the position `(ri, ci)` at which we should operate the `ith` operation.

Return _an array of integers_ `answer` _where_ `answer[i]` _is the number of islands after turning the cell_ `(ri, ci)` _into a land_.

An **island** is surrounded by water and is formed by connecting adjacent lands horizontally or vertically. You may assume all four edges of the grid are all surrounded by water.

**Example 1:**

```
Input:  m = 3, n = 3, positions = [[0,0],[0,1],[1,2],[2,1]]
Output: [1,1,2,3]

Explanation:
- Initially, the 2d grid is filled with water.
- Operation #1: addLand(0, 0) turns the water at grid[0][0] into a land. We have 1 island.
- Operation #2: addLand(0, 1) turns the water at grid[0][1] into a land. We still have 1 island.
- Operation #3: addLand(1, 2) turns the water at grid[1][2] into a land. We have 2 islands.
- Operation #4: addLand(2, 1) turns the water at grid[2][1] into a land. We have 3 islands.
```

**Example 2:**

```
Input:  m = 1, n = 1, positions = [[0,0]]
Output: [1]
```

**Constraints:**

- `1 <= m, n, positions.length <= 10^4`
- `1 <= m * n <= 10^4`
- `positions[i].length == 2`
- `0 <= ri < m`
- `0 <= ci < n`

**Follow up:** Could you solve it in time complexity `O(k log(mn))`, where `k == positions.length`?

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2024          |
| Amazon    | ★★★★☆ High       | 2023          |
| Facebook  | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Union-Find / Disjoint Set** — incremental dynamic connectivity as land is added; each merge decrements the island count → see [`/dsa/union_find.md`](/dsa/union_find.md)
- **Graph BFS/DFS (flood fill)** — the brute-force recount treats land as a graph and floods each component → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Matrix traversal** — mapping 2D cells to flat indices `r*n + c` → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS flood fill after each op | O(k·m·n) | O(m·n) | Small grids / correctness baseline; too slow at scale |
| 2 | Union-Find (Optimal) | O(k·α(m·n)) ≈ O(k) | O(m·n) | Dynamic connectivity — the intended solution, meets the follow-up |

---

## Approach 1 — BFS Flood Fill After Each Operation

### Intuition
The simplest correct method: maintain the actual grid, and after each add-land, recount the islands from scratch by scanning every cell and BFS-flooding each unvisited land component. Correct but O(m·n) per operation.

### Algorithm
1. Keep a boolean grid; for each position set that cell to land.
2. Recount: iterate all cells; each unvisited land cell begins a BFS that marks its whole island; increment a counter per island.
3. Append the count after each operation.

### Complexity
- **Time:** O(k·m·n) — k recounts, each a full grid scan.
- **Space:** O(m·n) — grid plus per-recount visited array.

### Code
```go
func floodFillCount(m, n int, positions [][]int) []int {
	grid := make([][]bool, m) // true = land
	for i := range grid {
		grid[i] = make([]bool, n)
	}
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // 4-neighbours
	result := make([]int, 0, len(positions))

	for _, p := range positions {
		grid[p[0]][p[1]] = true // add this land

		visited := make([][]bool, m)
		for i := range visited {
			visited[i] = make([]bool, n)
		}
		count := 0
		for r := 0; r < m; r++ {
			for c := 0; c < n; c++ {
				if grid[r][c] && !visited[r][c] { // new island root
					count++
					queue := [][2]int{{r, c}} // BFS frontier
					visited[r][c] = true
					for len(queue) > 0 {
						cell := queue[0]
						queue = queue[1:]
						for _, d := range dirs {
							nr, nc := cell[0]+d[0], cell[1]+d[1]
							if nr >= 0 && nr < m && nc >= 0 && nc < n &&
								grid[nr][nc] && !visited[nr][nc] {
								visited[nr][nc] = true
								queue = append(queue, [2]int{nr, nc})
							}
						}
					}
				}
			}
		}
		result = append(result, count)
	}
	return result
}
```

### Dry Run
`m=3, n=3, positions = [[0,0],[0,1],[1,2],[2,1]]`.

| Op | land added | grid land cells | islands after recount |
|---|---|---|---|
| 1 | (0,0) | {(0,0)} | 1 |
| 2 | (0,1) | {(0,0),(0,1)} adjacent | 1 |
| 3 | (1,2) | {(0,0),(0,1)} + {(1,2)} separate | 2 |
| 4 | (2,1) | above + {(2,1)} separate | 3 |

Result: `[1 1 2 3]`.

---

## Approach 2 — Union-Find / Disjoint Set (Optimal)

### Intuition
Adding one land cell provisionally creates a new island (`count++`). Then look at its up-to-four already-land neighbours: each one that belongs to a *different* island, when merged, removes one island (`count--`). Union-Find with path compression and union by rank makes each find/union nearly O(1), so we never rescan the grid. Cell `(r, c)` maps to flat index `r*n + c`.

### Algorithm
1. `parent[i] = i`; a separate `seen[]` marks which cells are land.
2. For each position: if it is already land, repeat the current count and skip.
3. Otherwise mark it land and `count++`.
4. For each in-bounds land neighbour, `union` it in; every successful merge does `count--`.
5. Append `count` after each operation.

### Complexity
- **Time:** O(k·α(m·n)) ≈ O(k) — near-constant find/union per operation (α is the inverse Ackermann function). Comfortably within the O(k·log(mn)) follow-up.
- **Space:** O(m·n) — parent, rank, and seen arrays.

### Code
```go
func unionFind(m, n int, positions [][]int) []int {
	total := m * n
	parent := make([]int, total) // disjoint-set parent pointers
	rank := make([]int, total)   // tree height hint for union by rank
	seen := make([]bool, total)  // whether a cell is land yet
	for i := range parent {
		parent[i] = i // each node is initially its own root
	}

	var find func(x int) int
	find = func(x int) int {
		for parent[x] != x { // walk up to the root
			parent[x] = parent[parent[x]] // path halving (compression)
			x = parent[x]
		}
		return x
	}
	union := func(a, b int) bool {
		ra, rb := find(a), find(b)
		if ra == rb {
			return false // already the same island — no merge
		}
		if rank[ra] < rank[rb] { // attach shorter tree under taller
			ra, rb = rb, ra
		}
		parent[rb] = ra
		if rank[ra] == rank[rb] {
			rank[ra]++ // heights tied → resulting tree grew by one
		}
		return true
	}

	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	result := make([]int, 0, len(positions))
	count := 0 // running number of islands

	for _, p := range positions {
		r, c := p[0], p[1]
		idx := r*n + c
		if seen[idx] { // duplicate add-land: island count unchanged
			result = append(result, count)
			continue
		}
		seen[idx] = true
		count++ // provisionally a brand-new island

		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr < 0 || nr >= m || nc < 0 || nc >= n {
				continue // off-grid neighbour
			}
			nidx := nr*n + nc
			if !seen[nidx] {
				continue // neighbour is still water
			}
			if union(idx, nidx) { // merged two distinct islands
				count-- // two islands became one
			}
		}
		result = append(result, count)
	}
	return result
}
```

### Dry Run
`m=3, n=3, positions = [[0,0],[0,1],[1,2],[2,1]]` (idx = r*3 + c).

| Op | cell (idx) | count after ++ | neighbour merges | count | result |
|---|---|---|---|---|---|
| 1 | (0,0)=0 | 1 | none land | 1 | 1 |
| 2 | (0,1)=1 | 2 | (0,0)=0 land → union ⇒ count-- | 1 | 1 |
| 3 | (1,2)=5 | 2 | (0,2),(2,2),(1,1) all water | 2 | 2 |
| 4 | (2,1)=7 | 3 | (1,1),(2,0),(2,2) all water | 3 | 3 |

Result: `[1 1 2 3]`.

---

## Key Takeaways

- **Dynamic connectivity ⇒ Union-Find.** When elements/edges are added incrementally and you must report connected-component counts along the way, DSU is the natural fit — no rescans.
- Maintain the component count **incrementally**: `+1` for each new node, `−1` for each successful merge.
- Flatten 2D coordinates to `r*n + c` to index a 1D DSU array.
- **Path compression + union by rank** give near-constant amortized operations (inverse Ackermann), meeting the O(k·log(mn)) follow-up with room to spare.
- Guard against **duplicate add-land** operations (same cell twice) so they do not spuriously bump the count.

---

## Related Problems

- LeetCode #200 — Number of Islands (static grid, one flood fill)
- LeetCode #547 — Number of Provinces (Union-Find on an adjacency matrix)
- LeetCode #684 — Redundant Connection (Union-Find cycle detection)
- LeetCode #261 — Graph Valid Tree (connectivity via Union-Find)
- LeetCode #1319 — Number of Operations to Make Network Connected (DSU components)
