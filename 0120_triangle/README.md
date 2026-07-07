# 0120 — Triangle

> LeetCode #120 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming

---

## Problem Statement

Given a `triangle` array, return the minimum path sum from top to bottom.

For each step, you may move to an adjacent number of the row below. More formally, if you are on index `i` on the current row, you may move to either index `i` or index `i + 1` on the next row.

**Example 1:**
```
Input: triangle = [[2],[3,4],[6,5,7],[4,1,8,3]]
Output: 11
Explanation: The path 2 → 3 → 5 → 1 has a sum of 11.
```

**Example 2:**
```
Input: triangle = [[-10]]
Output: -10
```

**Constraints:**
- `1 <= triangle.length <= 200`
- `triangle[i].length == i + 1`
- `-10^4 <= triangle[i][j] <= 10^4`

**Follow up:** Could you do this using only O(n) extra space, where n is the total number of rows in the triangle?

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Google    | ★★★★☆ High  | 2024          |
| Microsoft | ★★★☆☆ Medium | 2023          |
| Bloomberg | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **DP** — minimum-cost path on a DAG → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Rolling Array** — O(n) space by working bottom-up

---

## Approaches Overview

| # | Approach               | Time  | Space  | When to use              |
|---|------------------------|-------|--------|--------------------------|
| 1 | Top-Down DP + Memo     | O(n²) | O(n²)  | Intuitive                |
| 2 | Bottom-Up DP (Optimal) | O(n²) | O(n)   | Satisfies follow-up      |

---

## Approach 1 — Top-Down DP with Memoization

### Intuition
`dp(row, col)` = minimum path sum from `(row, col)` to the bottom.
At the last row: `dp(n-1, col) = triangle[n-1][col]`.
Otherwise: `dp(row, col) = triangle[row][col] + min(dp(row+1, col), dp(row+1, col+1))`.

Memoize to avoid recomputation.

### Complexity
- **Time:** O(n²) — each of the n(n+1)/2 cells computed once.
- **Space:** O(n²) — memo table.

### Code
```go
func minimumTotal(triangle [][]int) int {
    n := len(triangle)
    memo := make([][]int, n)
    for i := range memo {
        memo[i] = make([]int, len(triangle[i]))
        for j := range memo[i] { memo[i][j] = -1<<30 }
    }
    var dp func(row, col int) int
    dp = func(row, col int) int {
        if row == n-1 { return triangle[row][col] }
        if memo[row][col] != -1<<30 { return memo[row][col] }
        l, r := dp(row+1, col), dp(row+1, col+1)
        val := triangle[row][col]
        if l < r { val += l } else { val += r }
        memo[row][col] = val
        return val
    }
    return dp(0, 0)
}
```

### Dry Run
`[[2],[3,4],[6,5,7],[4,1,8,3]]`:

| dp call       | returns |
|---------------|---------|
| dp(3,0)=4     | 4       |
| dp(3,1)=1     | 1       |
| dp(3,2)=8     | 8       |
| dp(3,3)=3     | 3       |
| dp(2,0)=6+min(4,1)=7 | 7 |
| dp(2,1)=5+min(1,8)=6 | 6 |
| dp(2,2)=7+min(8,3)=10 | 10 |
| dp(1,0)=3+min(7,6)=9 | 9 |
| dp(1,1)=4+min(6,10)=10 | 10 |
| dp(0,0)=2+min(9,10)=11 | 11 |

---

## Approach 2 — Bottom-Up DP (Optimal)

### Intuition
Work from the bottom up. Initialize `dp` with the last row. For each row above, update: `dp[col] = triangle[row][col] + min(dp[col], dp[col+1])`. After processing all rows, `dp[0]` is the answer.

### Algorithm
1. `dp = copy of triangle[n-1]`.
2. For `row = n-2..0`:
   - For `col = 0..row`: `dp[col] = triangle[row][col] + min(dp[col], dp[col+1])`.
3. Return `dp[0]`.

### Complexity
- **Time:** O(n²)
- **Space:** O(n)

### Code
```go
func minimumTotalBottomUp(triangle [][]int) int {
    n := len(triangle)
    dp := make([]int, n); copy(dp, triangle[n-1])
    for row := n-2; row >= 0; row-- {
        for col := 0; col <= row; col++ {
            if dp[col] < dp[col+1] { dp[col] = triangle[row][col] + dp[col] } else { dp[col] = triangle[row][col] + dp[col+1] }
        }
    }
    return dp[0]
}
```

### Dry Run
`dp = [4,1,8,3]`.

| row=2 | col=0: dp[0]=6+min(4,1)=7 | col=1: dp[1]=5+min(1,8)=6 | col=2: dp[2]=7+min(8,3)=10 |
|-------|---------------------------|---------------------------|----------------------------|

`dp = [7,6,10,3]`.

| row=1 | col=0: dp[0]=3+min(7,6)=9 | col=1: dp[1]=4+min(6,10)=10 |
|-------|---------------------------|------------------------------|

`dp = [9,10,10,3]`.

| row=0 | col=0: dp[0]=2+min(9,10)=11 |
|-------|------------------------------|

Return 11 ✓

---

## Key Takeaways
- Bottom-up avoids stack overhead and naturally produces O(n) space.
- Adjacent move constraint: from `(row, col)` you can go to `(row+1, col)` or `(row+1, col+1)`.
- Bottom-up: `dp[col] = triangle[row][col] + min(dp[col], dp[col+1])` — `dp[col+1]` not yet overwritten when we process `col`.

---

## Related Problems
- LeetCode #64 — Minimum Path Sum (grid, not triangle)
- LeetCode #931 — Minimum Falling Path Sum
- LeetCode #1289 — Minimum Falling Path Sum II
