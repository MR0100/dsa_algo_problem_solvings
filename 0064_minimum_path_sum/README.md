# 0064 — Minimum Path Sum

> LeetCode #64 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, Matrix

---

## Problem Statement

Given a `m x n` grid filled with non-negative numbers, find a path from top left to bottom right, which **minimizes** the sum of all numbers along its path.

**Note:** You can only move either **down** or **right** at any point in time.

**Example 1**
```
Input:  grid = [[1,3,1],[1,5,1],[4,2,1]]
Output: 7
Explanation: Because the path 1 → 3 → 1 → 1 → 1 minimizes the sum.
```

**Example 2**
```
Input:  grid = [[1,2,3],[4,5,6]]
Output: 12
```

**Constraints**
- `m == grid.length`
- `n == grid[i].length`
- `1 <= m, n <= 200`
- `0 <= grid[i][j] <= 200`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★★ Very High | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming** — `dp[r][c]` = minimum cost to reach (r,c); choose the cheaper of "from above" and "from left."
- **In-Place DP** — reuse the input grid to avoid extra memory.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP 2D Table | O(m × n) | O(m × n) | Explicit; preserves input |
| 2 | DP In-Place ✅ | O(m × n) | O(1) | Optimal space; acceptable to modify input |
| 3 | DP Rolling Row | O(m × n) | O(n) | Space-optimized; preserves input |

---

## Approach 1 — DP 2D Table

### Intuition
`dp[r][c]` = minimum cost to reach cell (r,c).
- `dp[0][0] = grid[0][0]`
- First row: accumulate left-to-right (only right moves possible).
- First col: accumulate top-to-bottom (only down moves possible).
- Interior: `dp[r][c] = min(dp[r-1][c], dp[r][c-1]) + grid[r][c]`.

### Complexity
- **Time:** O(m × n).
- **Space:** O(m × n).

### Code
```go
// dpBottomUp solves Minimum Path Sum using a 2D DP table.
//
// Time:  O(m × n)
// Space: O(m × n)
func dpBottomUp(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	dp[0][0] = grid[0][0]
	for r := 1; r < m; r++ {
		dp[r][0] = dp[r-1][0] + grid[r][0]
	}
	for c := 1; c < n; c++ {
		dp[0][c] = dp[0][c-1] + grid[0][c]
	}
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			if dp[r-1][c] < dp[r][c-1] {
				dp[r][c] = dp[r-1][c] + grid[r][c]
			} else {
				dp[r][c] = dp[r][c-1] + grid[r][c]
			}
		}
	}
	return dp[m-1][n-1]
}
```

### Dry Run — `grid = [[1,3,1],[1,5,1],[4,2,1]]`

| step | cell | computation | dp |
|------|------|-------------|----|
| init | dp[0][0] | grid[0][0] | 1 |
| first col | dp[1][0] | dp[0][0]+grid[1][0]=1+1 | 2 |
| first col | dp[2][0] | dp[1][0]+grid[2][0]=2+4 | 6 |
| first row | dp[0][1] | dp[0][0]+grid[0][1]=1+3 | 4 |
| first row | dp[0][2] | dp[0][1]+grid[0][2]=4+1 | 5 |
| interior | dp[1][1] | min(dp[0][1],dp[1][0])+5=min(4,2)+5 | 7 |
| interior | dp[1][2] | min(dp[0][2],dp[1][1])+1=min(5,7)+1 | 6 |
| interior | dp[2][1] | min(dp[1][1],dp[2][0])+2=min(7,6)+2 | 8 |
| interior | dp[2][2] | min(dp[1][2],dp[2][1])+1=min(6,8)+1 | 7 |

Return `dp[2][2] = 7` ✓ (path: 1→3→1→1→1)

---

## Approach 2 — DP In-Place (Recommended ✅)

### Intuition
Same recurrence; write results back into `grid[r][c]` directly.

No extra allocation needed. If modifying input is acceptable (it usually is in interviews), this is cleanest.

### Complexity
- **Time:** O(m × n).
- **Space:** O(1).

### Code
```go
// dpInPlace solves Minimum Path Sum by reusing the input grid as the DP table.
//
// Time:  O(m × n)
// Space: O(1)  — no extra allocation (grid reused).
func dpInPlace(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	// first row accumulates left-to-right
	for c := 1; c < n; c++ {
		grid[0][c] += grid[0][c-1]
	}
	// first col accumulates top-to-bottom
	for r := 1; r < m; r++ {
		grid[r][0] += grid[r-1][0]
	}
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			if grid[r-1][c] < grid[r][c-1] {
				grid[r][c] += grid[r-1][c]
			} else {
				grid[r][c] += grid[r][c-1]
			}
		}
	}
	return grid[m-1][n-1]
}
```

### Dry Run — `grid = [[1,3,1],[1,5,1],[4,2,1]]`
```
First row accumulation: [1, 1+3=4, 4+1=5] → grid[0] = [1,4,5]
First col: [1, 1+1=2, 2+4=6] → grid[*][0] = [1,2,6]

Interior:
  grid[1][1] = min(grid[0][1], grid[1][0]) + 5 = min(4,2)+5 = 7
  grid[1][2] = min(grid[0][2], grid[1][1]) + 1 = min(5,7)+1 = 6
  grid[2][1] = min(grid[1][1], grid[2][0]) + 2 = min(7,6)+2 = 8
  grid[2][2] = min(grid[1][2], grid[2][1]) + 1 = min(6,8)+1 = 7

Return grid[2][2] = 7 ✓ (path: 1→3→1→1→1)
```

---

## Approach 3 — DP Rolling Row

### Intuition
Maintain one row. For each new row, update `dp[0] += grid[r][0]` (only from above), then `dp[c] = min(dp[c], dp[c-1]) + grid[r][c]`.

Before updating `dp[c]`: it holds the min-cost to reach the cell above (row r-1, col c). After updating `dp[c-1]`: it holds the cost of the left neighbor in the current row.

### Complexity
- **Time:** O(m × n).
- **Space:** O(n).

### Code
```go
// dpRolling solves Minimum Path Sum with O(n) space.
//
// Time:  O(m × n)
// Space: O(n)
func dpRolling(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	dp := make([]int, n)
	dp[0] = grid[0][0]
	for c := 1; c < n; c++ {
		dp[c] = dp[c-1] + grid[0][c] // first row init
	}
	for r := 1; r < m; r++ {
		dp[0] += grid[r][0] // first col: only can come from above
		for c := 1; c < n; c++ {
			top := dp[c] // before update = dp[r-1][c]
			left := dp[c-1] // after update = dp[r][c-1]
			if top < left {
				dp[c] = top + grid[r][c]
			} else {
				dp[c] = left + grid[r][c]
			}
		}
	}
	return dp[n-1]
}
```

### Dry Run — `grid = [[1,3,1],[1,5,1],[4,2,1]]`

| step | action | dp after |
|------|--------|----------|
| init row 0 | dp[0]=1; dp[1]=dp[0]+3=4; dp[2]=dp[1]+1=5 | [1,4,5] |
| r=1, first col | dp[0]+=grid[1][0]=1+1 | [2,4,5] |
| r=1, c=1 | dp[1]=min(top=4,left=2)+5=2+5 | [2,7,5] |
| r=1, c=2 | dp[2]=min(top=5,left=7)+1=5+1 | [2,7,6] |
| r=2, first col | dp[0]+=grid[2][0]=2+4 | [6,7,6] |
| r=2, c=1 | dp[1]=min(top=7,left=6)+2=6+2 | [6,8,6] |
| r=2, c=2 | dp[2]=min(top=6,left=8)+1=6+1 | [6,8,7] |

Return `dp[n-1] = dp[2] = 7` ✓

---

## Key Takeaways

- **Reuse input grid** — in interviews, asking "can I modify input?" before doing in-place DP is a professional signal; the answer is usually yes for this type of problem.
- **This is #62 with `min` instead of `+`** — counting paths vs. finding min cost uses the same grid DP skeleton. Switching `dp[r-1][c] + dp[r][c-1]` to `min(dp[r-1][c], dp[r][c-1]) + grid[r][c]` is the only change.
- **Rolling row works because we only look one row back** — the rolling pattern applies whenever `dp[r][c]` depends only on `dp[r-1][c]` and `dp[r][c-1]`.

---

## Related Problems

- LeetCode #62 — Unique Paths (count paths; same structure, different operation)
- LeetCode #63 — Unique Paths II (obstacles; same structure)
- LeetCode #120 — Triangle (minimum path sum in triangular grid)
- LeetCode #931 — Minimum Falling Path Sum (can move to adjacent columns too)
