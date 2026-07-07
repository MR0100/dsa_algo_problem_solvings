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

---

## Approach 2 — DP In-Place (Recommended ✅)

### Intuition
Same recurrence; write results back into `grid[r][c]` directly.

No extra allocation needed. If modifying input is acceptable (it usually is in interviews), this is cleanest.

### Complexity
- **Time:** O(m × n).
- **Space:** O(1).

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
