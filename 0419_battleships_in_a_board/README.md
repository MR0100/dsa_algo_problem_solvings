# 0419 — Battleships in a Board

> LeetCode #419 · Difficulty: Medium
> **Categories:** Array, Matrix, Depth-First Search

---

## Problem Statement

Given an `m x n` matrix `board` where each cell is a battleship `'X'` or empty `'.'`, return *the number of the battleships on* `board`.

Battleships can only be placed horizontally or vertically on `board`. In other words, they can only be made of the shape `1 x k` (`1` row, `k` columns) or `k x 1` (`k` rows, `1` column), where `k` can be of any size. At least one horizontal or vertical cell separates between two battleships (i.e., there are no adjacent battleships).

**Example 1:**

```
Input: board = [["X",".",".","X"],[".",".",".","X"],[".",".",".","X"]]
Output: 2
```

**Example 2:**

```
Input: board = [["."]]
Output: 0
```

**Constraints:**

- `m == board.length`
- `n == board[i].length`
- `1 <= m, n <= 200`
- `board[i][j]` is either `'.'` or `'X'`.

**Follow up:** Could you do it in one-pass, using only `O(1)` extra memory and without modifying the values `board`?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Matrix traversal** — the optimal solution is a single row-major scan of the grid with two constant-time look-back checks (cell above, cell to the left) → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Graph BFS/DFS (connected components)** — the baseline treats each ship as a connected blob of `'X'` cells and counts components via flood fill → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Flood Fill (connected components) | O(m·n) | O(m·n) | General "count blobs"; works even if ships weren't straight/separated |
| 2 | Count Ship Heads (Optimal) | O(m·n) | O(1) | The intended answer — one pass, no visited array, board untouched |

---

## Approach 1 — Brute Force (Flood Fill / DFS Connected Components)

### Intuition

Ignore the special guarantees and view the board as a graph: each `'X'` is a node, edges join 4-directionally adjacent `'X'` cells, and every battleship is one **connected component**. Counting ships is then counting components. Sweep the grid; the first time you touch an unvisited `'X'`, that's a brand-new ship, so increment and flood-fill its entire body (marking cells visited) to avoid recounting. This is more general than needed — it would still work if ships were L-shaped or touching — which is exactly why it's the natural first attempt.

### Algorithm

1. Allocate a `visited` grid.
2. For each cell in row-major order: if it's `'X'` and not visited, do `count++` and DFS-flood all connected `'X'` cells, marking them visited.
3. Return `count`.

### Complexity

- **Time:** O(m·n) — each cell is entered a constant number of times across all floods.
- **Space:** O(m·n) — the `visited` grid, plus recursion depth up to a full row or column.

### Code

```go
func floodFill(board [][]byte) int {
	m, n := len(board), len(board[0])
	visited := make([][]bool, m)
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // N, S, W, E

	var dfs func(r, c int)
	dfs = func(r, c int) {
		// Stop at walls, water, or already-counted cells.
		if r < 0 || r >= m || c < 0 || c >= n || board[r][c] != 'X' || visited[r][c] {
			return
		}
		visited[r][c] = true // absorb this cell into the current ship
		for _, d := range dirs {
			dfs(r+d[0], c+d[1]) // spread along the ship's body
		}
	}

	count := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if board[r][c] == 'X' && !visited[r][c] {
				count++   // first time we touch this ship
				dfs(r, c) // consume the rest of it
			}
		}
	}
	return count
}
```

### Dry Run

Example 1: board
```
X . . X
. . . X
. . . X
```

| Scan cell | value | visited? | Action | count | flood marks |
|-----------|-------|----------|--------|-------|-------------|
| (0,0) | X | no | new ship → DFS | 1 | (0,0) only (neighbours are `.`) |
| (0,1),(0,2) | . | — | skip | 1 | — |
| (0,3) | X | no | new ship → DFS down | 2 | (0,3),(1,3),(2,3) |
| (1,3),(2,3) | X | yes | already flooded → skip | 2 | — |
| rest | . | — | skip | 2 | — |

Result: `2` ✔ — one single-cell ship and one vertical 3-cell ship.

---

## Approach 2 — Count Ship "Heads" — One Pass, O(1) Space (Optimal)

### Intuition

Each battleship is a straight horizontal or vertical line, and no two ships touch. That means every ship has a **unique, identifiable cell**: its top-left end — the `'X'` that has no `'X'` directly above it and no `'X'` directly to its left. A vertical ship's head is its topmost cell; a horizontal ship's head is its leftmost cell; a single cell is its own head. Every other `'X'` is a *continuation* (it has an `'X'` neighbour above or to the left). So we don't need a visited array or flood fill at all: scan once and count only the heads, using two look-back comparisons per cell. Nothing is written to the board — this directly satisfies the follow-up (one pass, O(1) space, no mutation).

### Algorithm

1. Scan every cell `(r, c)` in row-major order.
2. Skip if `board[r][c] != 'X'`.
3. Skip if `r > 0 && board[r-1][c] == 'X'` (continuation of a vertical ship).
4. Skip if `c > 0 && board[r][c-1] == 'X'` (continuation of a horizontal ship).
5. Otherwise it's a fresh ship head → `count++`.
6. Return `count`.

### Complexity

- **Time:** O(m·n) — a single pass, constant work per cell.
- **Space:** O(1) — just the counter; the board is never modified.

### Code

```go
func countHeads(board [][]byte) int {
	m, n := len(board), len(board[0])
	count := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if board[r][c] != 'X' {
				continue // water — skip
			}
			// A ship's body extends right/down from its head, so a cell with an
			// 'X' above OR to the left is a continuation, not a new ship.
			if r > 0 && board[r-1][c] == 'X' {
				continue // part of a vertical ship already counted at its top
			}
			if c > 0 && board[r][c-1] == 'X' {
				continue // part of a horizontal ship already counted at its left
			}
			count++ // this 'X' is the unique top-left head of a new ship
		}
	}
	return count
}
```

### Dry Run

Example 1: board
```
X . . X
. . . X
. . . X
```

| Cell | value | up == 'X'? | left == 'X'? | Head? | count |
|------|-------|------------|--------------|-------|-------|
| (0,0) | X | — (r=0) | — (c=0) | yes | 1 |
| (0,3) | X | — (r=0) | (0,2)=`.` no | yes | 2 |
| (1,3) | X | (0,3)=X yes | — | no (skip) | 2 |
| (2,3) | X | (1,3)=X yes | — | no (skip) | 2 |
| all `.` | . | — | — | — | 2 |

Result: `2` ✔ — exactly the two heads `(0,0)` and `(0,3)` are counted; the ship body cells `(1,3),(2,3)` are correctly skipped.

---

## Key Takeaways

- **Give each object a canonical cell.** For shapes with a guaranteed orientation, count a unique representative (here the top-left head) instead of flood-filling the whole object. This turns O(m·n) space into O(1) and drops the visited array entirely.
- **Look-back beats mark-visited.** Checking the already-processed neighbours (up and left in a row-major scan) tells you whether the current cell continues a prior object — a common trick for one-pass, read-only grid counting.
- **The problem's constraints ARE the algorithm.** "Straight lines, always separated" is what makes the head unique; without those guarantees you'd fall back to connected-components (Approach 1). Always ask what the stated invariants let you skip.
- **Connected components is the general hammer;** reach for the O(1) head-count only when the shape guarantees permit it.

---

## Related Problems

- LeetCode #200 — Number of Islands (connected components, general blobs)
- LeetCode #695 — Max Area of Island (component sizes via flood fill)
- LeetCode #1254 — Number of Closed Islands (border-aware component counting)
- LeetCode #463 — Island Perimeter (per-cell edge counting, one pass)
