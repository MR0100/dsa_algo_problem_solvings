# 0062 — Unique Paths

> LeetCode #62 · Difficulty: Medium
> **Categories:** Math, Dynamic Programming, Combinatorics

---

## Problem Statement

There is a robot on an `m x n` grid. The robot is initially located at the **top-left corner** (i.e., `grid[0][0]`). The robot tries to move to the **bottom-right corner** (i.e., `grid[m - 1][n - 1]`). The robot can only move either **down** or **right** at any point in time.

Given the two integers `m` and `n`, return the number of possible unique paths that the robot can take to reach the bottom-right corner.

**Example 1**
```
Input:  m = 3, n = 7
Output: 28
```

**Example 2**
```
Input:  m = 3, n = 2
Output: 3
Explanation: From the top-left corner, there are a total of 3 ways:
1. Right -> Down -> Down
2. Down -> Down -> Right
3. Down -> Right -> Down
```

**Constraints**
- `1 <= m, n <= 100`
- The answer will be less than or equal to `2 × 10⁹`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★★☆ High      | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming** — `dp[r][c]` = paths to reach cell (r,c); builds on sub-problems from above and left.
- **Combinatorics** — closed-form: C(m+n-2, m-1) binomial coefficient.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Memoization (Top-Down DP) | O(m × n) | O(m × n) | Good learning step |
| 2 | DP 2D Table | O(m × n) | O(m × n) | Standard interview answer |
| 3 | DP Rolling Row | O(m × n) | O(n) | Space-optimized DP |
| 4 | Combinatorics ✅ | O(min(m,n)) | O(1) | Mathematically optimal |

---

## Approach 1 — Memoization (Top-Down DP)

### Intuition
`paths(r, c)` = paths from (r,c) to bottom-right = `paths(r+1,c) + paths(r,c+1)`. Base: any cell in the last row or last column has exactly 1 path (only one direction to go).

### Complexity
- **Time:** O(m × n).
- **Space:** O(m × n) — memo table + O(m+n) stack.

---

## Approach 2 — DP 2D Table

### Intuition
`dp[r][c]` = unique paths to (r,c). Initialize all border cells to 1 (only one way to reach them). Fill interior cells: `dp[r][c] = dp[r-1][c] + dp[r][c-1]`.

### Dry Run — `m=3, n=3`
```
dp[0][*] = [1,1,1]   (top row: all 1)
dp[*][0] = [1,1,1]   (left col: all 1)

dp[1][1] = dp[0][1] + dp[1][0] = 1+1 = 2
dp[1][2] = dp[0][2] + dp[1][1] = 1+2 = 3
dp[2][1] = dp[1][1] + dp[2][0] = 2+1 = 3
dp[2][2] = dp[1][2] + dp[2][1] = 3+3 = 6 ✓
```

### Complexity
- **Time:** O(m × n).
- **Space:** O(m × n).

---

## Approach 3 — DP Rolling Row

### Intuition
Only the previous row is needed to compute the current row. Reuse one array: `dp[c] += dp[c-1]` processes the transition in-place.

### Complexity
- **Time:** O(m × n).
- **Space:** O(n).

---

## Approach 4 — Combinatorics (Recommended ✅)

### Intuition
The robot always makes exactly `(m-1)` down moves and `(n-1)` right moves. The total number of moves is `m+n-2`. We choose which `m-1` of them are downs:

**C(m+n-2, m-1) = (m+n-2)! / ((m-1)! × (n-1)!)**

Compute iteratively (multiply-then-divide each step) to avoid large intermediate values.

### Algorithm
```
total = m+n-2; k = min(m-1, n-1)
result = 1
for i = 0 to k-1:
  result = result * (total - i) / (i + 1)
```

### Complexity
- **Time:** O(min(m, n)).
- **Space:** O(1).

### Code
```go
func combinatorics(m, n int) int {
    total, k := m+n-2, min(m-1, n-1)
    result := 1
    for i := 0; i < k; i++ { result = result * (total-i) / (i+1) }
    return result
}
```

### Dry Run — `m=3, n=7`
```
total=8, k=min(2,6)=2
i=0: result = 1 * 8 / 1 = 8
i=1: result = 8 * 7 / 2 = 28 ✓
```

---

## Key Takeaways

- **DP grid is Pascal's triangle** — `dp[r][c]` equals the binomial coefficient C(r+c, r), which is why both DP and combinatorics give the same answer.
- **Iterative C(n,k) avoids overflow** — divide at each step: `result = result * (total-i) / (i+1)`. The intermediate value is always an integer because C(n,k) builds one step at a time.
- **`dp[c] += dp[c-1]`** in the rolling-row is correct because after the update, `dp[c-1]` represents the left neighbor (current row) and `dp[c]` before update represents the top neighbor (previous row).

---

## Related Problems

- LeetCode #63 — Unique Paths II (with obstacles; DP only, no combinatorics)
- LeetCode #64 — Minimum Path Sum (find min cost path, same grid DP)
- LeetCode #120 — Triangle (similar DP on a triangle grid)
