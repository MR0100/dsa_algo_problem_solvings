# 0063 — Unique Paths II

> LeetCode #63 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, Matrix

---

## Problem Statement

You are given an `m x n` integer array `grid`. There is a robot initially located at the **top-left corner** (i.e., `grid[0][0]`). The robot tries to move to the **bottom-right corner** (i.e., `grid[m - 1][n - 1]`). The robot can only move either **down** or **right** at any point in time.

An obstacle and space are marked as `1` or `0` respectively in `grid`. A path that the robot takes cannot include **any** square that is an obstacle.

Return the number of possible unique paths that the robot can take to reach the bottom-right corner.

**Example 1**
```
Input:  obstacleGrid = [[0,0,0],[0,1,0],[0,0,0]]
Output: 2
Explanation: The robot has 2 paths: Right-Right-Down-Down, Down-Down-Right-Right.
```

**Example 2**
```
Input:  obstacleGrid = [[0,1],[0,0]]
Output: 1
```

**Constraints**
- `m == obstacleGrid.length`
- `n == obstacleGrid[i].length`
- `1 <= m, n <= 100`
- `obstacleGrid[i][j]` is `0` or `1`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming** — same as #62; set blocked cells to 0.
- **Grid DP with Obstacles** — obstacle propagation: once a row/column prefix is blocked, all subsequent cells in that row/column are 0.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP 2D Table ✅ | O(m × n) | O(m × n) | Explicit; easy to understand |
| 2 | DP Rolling Array | O(m × n) | O(n) | Space-optimized version |

---

## Approach 1 — DP 2D Table (Recommended ✅)

### Intuition
`dp[r][c]` = number of unique paths to (r,c). If `obstacleGrid[r][c] == 1`, set `dp[r][c] = 0`. Otherwise:
- Edge cells: 1 if no obstacle blocked the path to them; 0 if any obstacle in the same row (for top row) or column (for left column) blocked it.
- Interior cells: `dp[r][c] = dp[r-1][c] + dp[r][c-1]`.

### Algorithm
```
if obstacleGrid[0][0] == 1: return 0
dp[0][0] = 1
first column: dp[r][0] = dp[r-1][0] if no obstacle else 0
first row: dp[0][c] = dp[0][c-1] if no obstacle else 0
interior: dp[r][c] = dp[r-1][c] + dp[r][c-1] if no obstacle else 0
return dp[m-1][n-1]
```

### Complexity
- **Time:** O(m × n).
- **Space:** O(m × n).

### Dry Run — `obstacleGrid = [[0,0,0],[0,1,0],[0,0,0]]`
```
dp[0] = [1,1,1]  (first row: no obstacles)
dp[1][0] = dp[0][0] = 1
dp[1][1] = 0     (obstacle!)
dp[1][2] = dp[0][2] + dp[1][1] = 1+0 = 1
dp[2][0] = dp[1][0] = 1
dp[2][1] = dp[1][1] + dp[2][0] = 0+1 = 1
dp[2][2] = dp[1][2] + dp[2][1] = 1+1 = 2 ✓
```

---

## Approach 2 — DP Rolling Array

### Intuition
Same recurrence with one array. When an obstacle is encountered, set `dp[c] = 0` to propagate the blockage. This correctly zeros out any paths that tried to pass through the obstacle.

### Complexity
- **Time:** O(m × n).
- **Space:** O(n).

---

## Key Takeaways

- **Obstacle in top-left → return 0 immediately** — no paths exist.
- **Obstacle propagation in first row/column** — once an obstacle is hit in the first row (or column), all cells to its right (or below) must be 0. This is natural with `dp[c] = dp[c-1]` (which propagates 0 forward).
- **Exact same structure as #62** — the only addition is the `if obstacle: dp[r][c] = 0` guard.
- **No combinatorics formula** — obstacles break the symmetry required for C(m+n-2, m-1); DP is the only approach.

---

## Related Problems

- LeetCode #62 — Unique Paths (no obstacles; combinatorics formula works)
- LeetCode #64 — Minimum Path Sum (minimize cost instead of counting paths)
- LeetCode #980 — Unique Paths III (must visit every non-obstacle cell; backtracking)
