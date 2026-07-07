# 0221 — Maximal Square

> LeetCode #221 · Difficulty: Medium
> **Categories:** Dynamic Programming, Matrix

---

## Problem Statement

Given an `m x n` binary `matrix` filled with `0`'s and `1`'s, find the largest square containing only `1`'s and return *its area*.

**Example 1:**

```
Input: matrix = [["1","0","1","0","0"],
                 ["1","0","1","1","1"],
                 ["1","1","1","1","1"],
                 ["1","0","0","1","0"]]
Output: 4
```

Explanation: The largest all-ones square has side 2, formed by the block of `1`'s in rows 1–2, columns 2–3. Its area is 2 × 2 = 4.

**Example 2:**

```
Input: matrix = [["0","1"],["1","0"]]
Output: 1
```

Explanation: No 2×2 all-ones block exists; the best is a single `1`, area 1.

**Example 3:**

```
Input: matrix = [["0"]]
Output: 0
```

Explanation: There are no `1`'s, so the largest square has area 0.

**Constraints:**

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 300`
- `matrix[i][j]` is `'0'` or `'1'`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **2D Dynamic Programming** — `dp[r][c]` = side of the largest all-ones square ending at `(r,c)`, built from three overlapping subproblems → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Rolling-Array Space Optimization** — the recurrence touches only the current and previous rows, so a single 1D row plus one scalar suffices → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Matrix Traversal** — row-major sweep over an `m x n` grid → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Expand Every Square) | O(m·n·min(m,n)²) | O(1) | Baseline / intuition; too slow for 300×300 worst cases |
| 2 | DP (2D Table) | O(m·n) | O(m·n) | Clearest correct DP; when you can afford the table |
| 3 | DP (1D Rolling Row) (Optimal) | O(m·n) | O(n) | Same speed, minimal memory — the answer to give |

---

## Approach 1 — Brute Force (Expand Every Square)

### Intuition
Every all-ones square has a top-left corner. Anchor a square at each `1` cell and grow it outward one ring at a time; the moment the new border contains a `0` (or runs off the grid), stop. The largest side that succeeded, over all anchors, gives the area.

### Algorithm
1. Initialise `best = 0`.
2. For each cell `(r,c)` equal to `'1'`:
   1. Start `side = 1`.
   2. While `(r+side, c+side)` stays in bounds and the L-shaped border added by growing to `side+1` (its new bottom row and new right column) is all `'1'`, increment `side`.
   3. Update `best = max(best, side)`.
3. Return `best * best`.

### Complexity
- **Time:** O(m·n·min(m,n)²) — each of the `m·n` anchors may scan borders totalling O(side²) work, with `side` up to `min(m,n)`.
- **Space:** O(1) — only scalar counters.

### Code
```go
func bruteForce(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0 // empty grid → no square
	}
	m, n := len(matrix), len(matrix[0])
	best := 0 // largest side length found so far
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if matrix[r][c] != '1' {
				continue // a square's top-left must itself be '1'
			}
			side := 1 // a lone '1' is already a 1×1 square
			for r+side < m && c+side < n && valid(matrix, r, c, side) {
				side++ // the (side+1)×(side+1) square is all ones too
			}
			if side > best {
				best = side // remember the biggest square so far
			}
		}
	}
	return best * best // area = side²
}

func valid(matrix [][]byte, r, c, side int) bool {
	for i := 0; i <= side; i++ {
		if matrix[r+side][c+i] != '1' { // new bottom row
			return false
		}
		if matrix[r+i][c+side] != '1' { // new right column
			return false
		}
	}
	return true
}
```

### Dry Run
Example 1, focusing on the winning anchor `(1,2)` (row 1, col 2):

| Step | Anchor (r,c) | side tried | Border all ones? | side after |
|------|--------------|------------|------------------|------------|
| 1 | (1,2) | grow to 2 → checks row 2 cols 2–3 + col 3 rows 1–2 | yes (all `1`) | 2 |
| 2 | (1,2) | grow to 3 → needs (4,·) | out of bounds (m=4) | stop at 2 |
| 3 | other anchors | — | none reaches side 3 | best stays 2 |

`best = 2` → area `2*2 = 4`. ✔

---

## Approach 2 — Dynamic Programming (2D Table)

### Intuition
Let `dp[r][c]` be the side of the largest all-ones square whose **bottom-right** corner is `(r,c)`. To have a square of side `k` here, you need squares of side `k-1` ending at the cell above, the cell to the left, and the cell diagonally up-left — the smallest of those three caps how far `(r,c)` can extend, plus 1 for the current cell.

### Algorithm
1. Allocate `dp` the same shape as `matrix`, all zeros.
2. Sweep row by row, column by column:
   - `'0'` cell → `dp[r][c] = 0`.
   - `'1'` on the top row or left column → `dp[r][c] = 1`.
   - otherwise `dp[r][c] = 1 + min(dp[r-1][c], dp[r][c-1], dp[r-1][c-1])`.
3. Track the maximum `dp` value; return its square.

### Complexity
- **Time:** O(m·n) — one O(1) recurrence per cell.
- **Space:** O(m·n) — the full DP table.

### Code
```go
func dp2D(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])
	dp := make([][]int, m) // dp[r][c] = largest square side ending at (r,c)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	best := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if matrix[r][c] != '1' {
				continue // '0' cell ends no square → dp stays 0
			}
			if r == 0 || c == 0 {
				dp[r][c] = 1 // border cell: at most a 1×1 square
			} else {
				dp[r][c] = 1 + min3(dp[r-1][c], dp[r][c-1], dp[r-1][c-1])
			}
			if dp[r][c] > best {
				best = dp[r][c] // track the biggest square side
			}
		}
	}
	return best * best
}
```

### Dry Run
Example 1 — `dp` table (row-major). Rows/cols 0-indexed; `.` = 0:

| r\c | 0 | 1 | 2 | 3 | 4 |
|-----|---|---|---|---|---|
| 0 | 1 | . | 1 | . | . |
| 1 | 1 | . | 1 | 1 | 1 |
| 2 | 1 | 1 | 1 | **2** | **2** |
| 3 | 1 | . | . | 1 | . |

At `(2,3)`: cell is `'1'`, neighbours `dp[1][3]=1`, `dp[2][2]=1`, `dp[1][2]=1` → `1 + min(1,1,1) = 2`. At `(2,4)`: neighbours `dp[1][4]=1`, `dp[2][3]=2`, `dp[1][3]=1` → `1 + min = 2`. Max `dp = 2` → area `4`. ✔

---

## Approach 3 — Dynamic Programming (1D Rolling Row) (Optimal)

### Intuition
The 2D recurrence reads only the current row and the one above it, so keeping the whole table wastes memory. Overwrite a single 1D array left-to-right; the only extra value you need is the *old* `dp[r-1][c-1]` (the diagonal), which you stash in `prev` right before overwriting the cell that will become the next column's top-left.

### Algorithm
1. Keep `dp` of length `n+1` (index shifted by 1 so column 0 has a free zero to its left).
2. For each row, reset `prev = 0`, then for each column `c`:
   1. `temp = dp[c+1]` (save the old top value = next iteration's top-left).
   2. `'1'` → `dp[c+1] = 1 + min(dp[c+1] (top), dp[c] (left), prev (top-left))`.
   3. `'0'` → `dp[c+1] = 0`.
   4. `prev = temp`.
3. Track the max; return its square.

### Complexity
- **Time:** O(m·n) — single pass, O(1) per cell.
- **Space:** O(n) — one rolling row plus two scalars.

### Code
```go
func dp1D(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])
	dp := make([]int, n+1) // dp[c+1] = square side ending at current row, column c
	best := 0
	for r := 0; r < m; r++ {
		prev := 0 // dp[r-1][c-1]: top-left diagonal, starts 0 each row
		for c := 0; c < n; c++ {
			temp := dp[c+1] // save dp[r-1][c] before overwriting (becomes next prev)
			if matrix[r][c] == '1' {
				dp[c+1] = 1 + min3(dp[c+1], dp[c], prev)
				if dp[c+1] > best {
					best = dp[c+1]
				}
			} else {
				dp[c+1] = 0 // '0' cell resets the running square
			}
			prev = temp // this row's dp[c] is next column's top-left
		}
	}
	return best * best
}
```

### Dry Run
Example 1, tracing row 2 (the row that produces the winning `2`). Entering row 2, `dp` (indices 1..5) holds the row-1 sides: `[_, 1,0,1,1,1]`. `prev` resets to 0.

| c | matrix[2][c] | temp=dp[c+1] (top) | dp[c] (left) | prev (top-left) | new dp[c+1] | best |
|---|--------------|--------------------|--------------|-----------------|-------------|------|
| 0 | 1 | 1 | dp[0]=0 | 0 | 1 | 1 |
| 1 | 1 | 0 | 1 | 1 | 1 | 1 |
| 2 | 1 | 1 | 1 | 0 | 1 | 1 |
| 3 | 1 | 1 | 1 | 1 | **2** | **2** |
| 4 | 1 | 1 | 2 | 1 | **2** | 2 |

`best = 2` → area `4`. ✔

---

## Key Takeaways
- The square-DP recurrence `1 + min(top, left, top-left)` is a reusable pattern: the *minimum* of three neighbours is what limits a growing square (contrast with rectangle problems, which need histogram/stack techniques).
- Any DP whose recurrence reads only the previous row collapses to O(n) space with a rolling array plus a saved diagonal scalar.
- Answer is the **area** (`side²`), a classic off-by-one trap — the DP naturally tracks the *side*.

---

## Related Problems
- LeetCode #1277 — Count Square Submatrices with All Ones (same DP, sum the `dp` values instead of taking the max)
- LeetCode #85 — Maximal Rectangle (rectangle variant; needs histogram + monotonic stack)
- LeetCode #84 — Largest Rectangle in Histogram (building block for #85)
- LeetCode #764 — Largest Plus Sign (directional DP over a grid)
