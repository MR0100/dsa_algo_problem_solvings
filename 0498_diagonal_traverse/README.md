# 0498 — Diagonal Traverse

> LeetCode #498 · Difficulty: Medium
> **Categories:** Array, Matrix, Simulation

---

## Problem Statement

Given an `m x n` matrix `mat`, return *an array of all the elements of the array in a diagonal order*.

**Example 1:**

```
Input: mat = [[1,2,3],[4,5,6],[7,8,9]]
Output: [1,2,4,7,5,3,6,8,9]
```

*(The traversal snakes along anti-diagonals: the first diagonal goes up-right, the next down-left, alternating.)*

**Example 2:**

```
Input: mat = [[1,2],[3,4]]
Output: [1,2,3,4]
```

**Constraints:**

- `m == mat.length`
- `n == mat[i].length`
- `1 <= m, n <= 10^4`
- `1 <= m * n <= 10^4`
- `-10^5 <= mat[i][j] <= 10^5`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Matrix Traversal** — the core skill is visiting grid cells in a non-row-major order; here the *anti-diagonal* order with alternating direction, driven either by the `r+c` invariant or by an edge-bouncing walk → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Group by Diagonal, Reverse Alternate | O(m·n) | O(m·n) | Clearest to reason about; when extra buffer is fine |
| 2 | Simulation with Direction Flips (Optimal) | O(m·n) | O(1) | Interview-preferred; O(1) extra space, single walk |

---

## Approach 1 — Group by Diagonal, Reverse Alternate

### Intuition

Every cell on the same anti-diagonal shares the sum `r + c`. There are `m + n − 1` anti-diagonals, indexed `d = 0 … m+n−2`. Bucket each cell under its `d`, appending in increasing row order (so each bucket is naturally top-to-bottom). The zig-zag output is then just: emit buckets in increasing `d`, but for **even** `d` read the bucket **reversed** (that diagonal travels up-right, i.e. bottom cell first), and for **odd** `d` read it as-is (down-left). The direction alternation falls straight out of the parity of `d`.

### Algorithm

1. Handle empty input → return `[]`.
2. Create `m+n−1` empty buckets.
3. For each `(r,c)`: append `mat[r][c]` to bucket `d = r+c` (rows scanned in increasing order).
4. For `d = 0 … m+n−2`: if `d` is even, append the bucket **reversed**; else append it **forward**.
5. Return the concatenation.

### Complexity

- **Time:** O(m·n) — each cell is bucketed once and later emitted once.
- **Space:** O(m·n) — the buckets store every cell before output (the output slice itself is required either way).

### Code

```go
func groupByDiagonal(mat [][]int) []int {
	if len(mat) == 0 || len(mat[0]) == 0 {
		return []int{}
	}
	m, n := len(mat), len(mat[0])
	diagonals := make([][]int, m+n-1) // one bucket per anti-diagonal d = r+c
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			d := r + c // anti-diagonal index; appended in increasing r (top→bottom)
			diagonals[d] = append(diagonals[d], mat[r][c])
		}
	}

	result := make([]int, 0, m*n)
	for d := 0; d < len(diagonals); d++ {
		bucket := diagonals[d]
		if d%2 == 0 {
			// even diagonal → travels UP-right, i.e. bottom row first: reverse it
			for i := len(bucket) - 1; i >= 0; i-- {
				result = append(result, bucket[i])
			}
		} else {
			// odd diagonal → travels DOWN-left, natural top→bottom order
			result = append(result, bucket...)
		}
	}
	return result
}
```

### Dry Run

Example 1: `mat = [[1,2,3],[4,5,6],[7,8,9]]`, so `m=n=3`, diagonals `d=0..4`.

Bucketing (each cell to `d = r+c`, rows increasing):

| d | cells appended (top→bottom) |
|---|-----------------------------|
| 0 | [1] |
| 1 | [2, 4] |
| 2 | [3, 5, 7] |
| 3 | [6, 8] |
| 4 | [9] |

Emit with parity rule:

| d | parity | action | contributes |
|---|--------|--------|-------------|
| 0 | even | reverse [1] | 1 |
| 1 | odd | forward [2,4] | 2, 4 |
| 2 | even | reverse [3,5,7] → [7,5,3] | 7, 5, 3 |
| 3 | odd | forward [6,8] | 6, 8 |
| 4 | even | reverse [9] | 9 |

Result: `[1, 2, 4, 7, 5, 3, 6, 8, 9]` ✔

---

## Approach 2 — Simulation with Direction Flips (Optimal)

### Intuition

Walk the grid cell-by-cell holding a `(row, col)` and a `direction`. Going **up-right** means `r--, c++`; going **down-left** means `r++, c--`. When a step would leave the grid, "bounce" to the start of the next diagonal and flip direction. The subtlety is the corner cells: when moving up and you are simultaneously in the top row *and* the right column, the **right-column** rule must win (drop down a row, don't step right), otherwise you skip a cell. Symmetrically, moving down, the **bottom-row** rule wins over the left column. Only O(1) extra state is needed.

### Algorithm

Repeat until `m·n` cells emitted; each step first records `mat[r][c]`, then moves:

**Moving up-right (`direction=+1`):**
1. If `c == n−1` (right edge): `r++`, flip to down-left.
2. Else if `r == 0` (top edge): `c++`, flip to down-left.
3. Else: `r--`, `c++`.

**Moving down-left (`direction=−1`):**
1. If `r == m−1` (bottom edge): `c++`, flip to up-right.
2. Else if `c == 0` (left edge): `r++`, flip to up-right.
3. Else: `r++`, `c--`.

### Complexity

- **Time:** O(m·n) — exactly one visit per cell.
- **Space:** O(1) — `r`, `c`, `direction` counters beyond the required output.

### Code

```go
func simulateWalk(mat [][]int) []int {
	if len(mat) == 0 || len(mat[0]) == 0 {
		return []int{}
	}
	m, n := len(mat), len(mat[0])
	result := make([]int, 0, m*n)
	r, c := 0, 0     // current cell, start at top-left
	direction := 1   // +1 = moving up-right, -1 = moving down-left

	for len(result) < m*n { // exactly m*n cells to emit
		result = append(result, mat[r][c]) // record the current cell

		if direction == 1 { // moving up-right (r--, c++)
			switch {
			case c == n-1: // hit right wall → drop down one row, flip to down-left
				r++
				direction = -1
			case r == 0: // hit top wall (and not right wall) → step right, flip
				c++
				direction = -1
			default: // free to keep moving up-right
				r--
				c++
			}
		} else { // moving down-left (r++, c--)
			switch {
			case r == m-1: // hit bottom wall → step right one column, flip to up-right
				c++
				direction = 1
			case c == 0: // hit left wall (and not bottom wall) → drop down, flip
				r++
				direction = 1
			default: // free to keep moving down-left
				r++
				c--
			}
		}
	}
	return result
}
```

### Dry Run

Example 1: `mat = [[1,2,3],[4,5,6],[7,8,9]]`, `m=n=3`. Start `(0,0)`, `direction=+1`.

| emit | cell (r,c) | value | direction now | edge hit → move | next (r,c), dir |
|------|-----------|-------|---------------|-----------------|-----------------|
| 1 | (0,0) | 1 | up | r==0 top → c++, flip | (0,1), down |
| 2 | (0,1) | 2 | down | free → r++,c-- | (1,0), down |
| 3 | (1,0) | 4 | down | c==0 left → r++, flip | (2,0), up |
| 4 | (2,0) | 7 | up | free → r--,c++ | (1,1), up |
| 5 | (1,1) | 5 | up | free → r--,c++ | (0,2), up |
| 6 | (0,2) | 3 | up | c==n-1 right → r++, flip | (1,2), down |
| 7 | (1,2) | 6 | down | free → r++,c-- | (2,1), down |
| 8 | (2,1) | 8 | down | r==m-1 bottom → c++, flip | (2,2), up |
| 9 | (2,2) | 9 | up | (loop ends, 9 cells) | — |

Result: `[1, 2, 4, 7, 5, 3, 6, 8, 9]` ✔

---

## Key Takeaways

- **Anti-diagonals are levels of `r + c`** (main diagonals are levels of `r − c`). Recognising this invariant instantly organises any diagonal problem.
- **Parity of the diagonal index controls direction:** even diagonals go one way, odd the other. Reversing alternate buckets is the simplest correct implementation.
- The **O(1) simulation** trades a buffer for careful edge handling. The gotcha is corner precedence: on the way *up*, test the **right column before the top row**; on the way *down*, test the **bottom row before the left column** — otherwise a corner cell is emitted twice or skipped.
- When two clean solutions share the same O(m·n) time, prefer the **O(1)-space** one in interviews, but keep the bucket version ready — it is far easier to prove correct under pressure.

---

## Related Problems

- LeetCode #54 — Spiral Matrix (another prescribed-order matrix walk)
- LeetCode #59 — Spiral Matrix II (generate in spiral order)
- LeetCode #1424 — Diagonal Traverse II (jagged rows; bucket by `r+c` with a heap/queue)
- LeetCode #48 — Rotate Image (index-mapping over a matrix)
- LeetCode #766 — Toeplitz Matrix (cells share `r−c` diagonals)
