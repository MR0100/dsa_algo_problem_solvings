# 0361 — Bomb Enemy

> LeetCode #361 · Difficulty: Medium
> **Categories:** Array, Matrix, Dynamic Programming, Prefix Sum

---

## Problem Statement

Given an `m x n` matrix `grid` where each cell is either a wall `'W'`, an enemy `'E'` or empty `'0'`, return *the maximum enemies you can kill using one bomb*. You can only place the bomb in an empty cell.

The bomb kills all the enemies in the same row and column from the planted point until it hits the wall since it is too strong to be destroyed.

**Example 1:**

```
Input: grid = [["0","E","0","0"],["E","0","W","E"],["0","E","0","0"]]
Output: 3
```

**Example 2:**

```
Input: grid = [["W","W","W"],["0","0","0"],["E","E","E"]]
Output: 1
```

**Constraints:**

- `m == grid.length`
- `n == grid[i].length`
- `1 <= m, n <= 500`
- `grid[i][j]` is either `'W'`, `'E'`, or `'0'`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Matrix Traversal** — the entire problem is a grid sweep; the bomb's blast propagates along a row and a column until a wall → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Prefix Sum / Running Count** — the optimal trick is to reuse an accumulated enemy count per wall-bounded row/column segment so each cell is O(1) → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Dynamic Programming (on the grid)** — `colHits[c]` carries a partial result forward across rows, recomputed only when a wall breaks the segment → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(m·n·(m+n)) | O(1) | Small grids, or to establish correctness quickly |
| 2 | Row/Column Running Count (Optimal) | O(m·n) | O(n) | The intended answer; each cell resolved in O(1) |

---

## Approach 1 — Brute Force

### Intuition

The statement is a simulation. Drop the bomb on an empty cell and it kills every enemy in the same row and column, stopping at the first wall in each of the four directions. So visit every empty cell, walk outward in all four directions, count enemies, and keep the best total.

### Algorithm

1. For each cell `(r,c)` equal to `'0'`:
   1. Walk **up** from `(r,c)`: count `'E'`, stop at `'W'` or the top border.
   2. Walk **down**, **left**, **right** the same way.
   3. Sum the four counts into `kills`.
2. Track the maximum `kills` seen across all empty cells.

### Complexity

- **Time:** O(m·n·(m+n)) — up to `m·n` empty cells, each scanning up to `m+n` cells across its row and column.
- **Space:** O(1) — a handful of counters.

### Code

```go
func bruteForce(grid [][]byte) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0 // no cells → nothing to bomb
	}
	m, n := len(grid), len(grid[0])
	best := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] != '0' {
				continue // bomb can only be placed on an empty cell
			}
			kills := 0
			// Walk up: decrement row until wall or top border.
			for i := r - 1; i >= 0 && grid[i][c] != 'W'; i-- {
				if grid[i][c] == 'E' {
					kills++ // an enemy in the blast path
				}
			}
			// Walk down.
			for i := r + 1; i < m && grid[i][c] != 'W'; i++ {
				if grid[i][c] == 'E' {
					kills++
				}
			}
			// Walk left.
			for j := c - 1; j >= 0 && grid[r][j] != 'W'; j-- {
				if grid[r][j] == 'E' {
					kills++
				}
			}
			// Walk right.
			for j := c + 1; j < n && grid[r][j] != 'W'; j++ {
				if grid[r][j] == 'E' {
					kills++
				}
			}
			if kills > best {
				best = kills // remember the strongest placement
			}
		}
	}
	return best
}
```

### Dry Run

Example 1, `grid = [["0","E","0","0"],["E","0","W","E"],["0","E","0","0"]]`. Best cell is `(1,1)` (empty).

| Direction from (1,1) | Path cells | Stops at | Enemies counted |
|----------------------|-----------|----------|-----------------|
| Up    | (0,1)=`E`          | top border | 1 |
| Down  | (2,1)=`E`          | bottom border | 1 |
| Left  | (1,0)=`E`          | left border | 1 |
| Right | (1,2)=`W`          | wall immediately | 0 |

`kills = 1+1+1+0 = 3`. No other empty cell beats it → answer `3` ✔

---

## Approach 2 — Row/Column Running Count (Optimal)

### Intuition

For a bomb placed anywhere inside one wall-bounded row segment, the number of enemies killed *along that row* is identical — it is just the enemy count of the segment. The count only changes when we cross a wall. So keep one `rowHits` accumulator that is recomputed at the start of each row segment (right after a wall or the left border) and reused for every empty cell in the segment. Do the same per column with `colHits[c]`, recomputed whenever a wall above breaks the column segment. At an empty cell the answer is simply `rowHits + colHits[c]`, computed in O(1).

### Algorithm

1. Maintain `rowHits` (enemies in the current row segment) and an array `colHits[c]` (enemies in column `c`'s current segment).
2. Sweep row by row, left to right. At `(r,c)`:
   1. If `c == 0` or the cell to the left is `'W'`, recount enemies rightward until the next wall → `rowHits`.
   2. If `r == 0` or the cell above is `'W'`, recount enemies downward until the next wall → `colHits[c]`.
   3. If `(r,c)` is `'0'`, candidate `= rowHits + colHits[c]`; update the max.
3. Return the max.

### Complexity

- **Time:** O(m·n) — each segment is recounted once, and its cost is charged across the cells of that segment; amortised O(1) per cell.
- **Space:** O(n) — the `colHits` array of running column counts.

### Code

```go
func runningCount(grid [][]byte) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])
	best := 0
	rowHits := 0                 // enemies in the current row segment
	colHits := make([]int, n)    // enemies in each column's current segment
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			// Recompute rowHits at the start of a new row segment: either the
			// left border, or just after a wall on the left.
			if c == 0 || grid[r][c-1] == 'W' {
				rowHits = 0
				for k := c; k < n && grid[r][k] != 'W'; k++ {
					if grid[r][k] == 'E' {
						rowHits++ // enemies reachable rightward before a wall
					}
				}
			}
			// Recompute colHits[c] at the start of a new column segment: either
			// the top border, or just after a wall above.
			if r == 0 || grid[r-1][c] == 'W' {
				colHits[c] = 0
				for k := r; k < m && grid[k][c] != 'W'; k++ {
					if grid[k][c] == 'E' {
						colHits[c]++ // enemies reachable downward before a wall
					}
				}
			}
			// Only empty cells can hold the bomb; sum row + column kills.
			if grid[r][c] == '0' {
				if total := rowHits + colHits[c]; total > best {
					best = total
				}
			}
		}
	}
	return best
}
```

### Dry Run

Example 1, `grid = [["0","E","0","0"],["E","0","W","E"],["0","E","0","0"]]`. Focus on the winning cell `(1,1)`.

| At cell | New row segment? | rowHits | New col segment? | colHits[c] | empty? | candidate |
|---------|------------------|---------|------------------|------------|--------|-----------|
| (0,0) `0` | yes (c=0): row 0 up to end has 1 `E` | 1 | yes (r=0): col 0 down = `0,E,0` → 1 `E` | 1 | yes | 1+1 = 2 |
| (1,0) `E` | yes (c=0): row 1 up to wall at col 2 → `E,0` → 1 `E` | 1 | reuse colHits[0]=1 | 1 | no | — |
| (1,1) `0` | reuse rowHits=1 (same segment, no wall left) | 1 | yes (above (0,1)=`E`? no wall) → reuse... above is `E` not `W`, so reuse colHits[1] | col 1 = `E,0,E` → 2 | 2 | yes | 1+2 = **3** |

`(1,1)` yields `3`, the maximum → answer `3` ✔

*(colHits[1] was computed at row 0 as enemies down column 1 with no walls: `E,0,E` = 2; it is reused at row 1 because the cell above `(1,1)` is `'0'`/`'E'`, not `'W'`.)*

---

## Key Takeaways

- **Reuse work across a wall-bounded segment.** Any bomb inside the same segment kills the same enemies along that axis — compute the count once per segment, not once per cell.
- **A row accumulator + a column-array accumulator** turns an O(m·n·(m+n)) simulation into O(m·n). The column state must persist across rows (an array), while the row state is a single scalar reset each row.
- **Recompute only at segment boundaries** — detect them cheaply by checking the neighbour to the left (`grid[r][c-1]=='W'`) or above (`grid[r-1][c]=='W'`).
- This "prefix count that resets on a delimiter" pattern recurs in grid problems (largest rectangle, gas station rings, matrix range sums with obstacles).

---

## Related Problems

- LeetCode #221 — Maximal Square (grid DP carrying partial results forward)
- LeetCode #304 — Range Sum Query 2D Immutable (2D prefix sums)
- LeetCode #85 — Maximal Rectangle (per-column running heights)
- LeetCode #073 — Set Matrix Zeroes (row/column propagation over a grid)
