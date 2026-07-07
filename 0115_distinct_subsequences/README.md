# 0115 — Distinct Subsequences

> LeetCode #115 · Difficulty: Hard
> **Categories:** String, Dynamic Programming

---

## Problem Statement

Given two strings `s` and `t`, return the number of distinct subsequences of `s` which equals `t`.

The test cases are generated so that the answer fits in a 32-bit signed integer.

**Example 1:**
```
Input: s = "rabbbit", t = "rabbit"
Output: 3
Explanation:
rabbbit → rabbit (choosing b at positions 2,3,4 — pick any two of three b's: C(3,2)=3 ways)
```

**Example 2:**
```
Input: s = "babgbag", t = "bag"
Output: 5
```

**Constraints:**
- `1 <= s.length, t.length <= 1000`
- `s` and `t` consist of English lowercase letters.

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Google    | ★★★★☆ High  | 2024          |
| Amazon    | ★★★★☆ High  | 2024          |
| Facebook  | ★★★☆☆ Medium | 2023          |
| Bloomberg | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **2D DP on strings** — `dp[i][j]` = count of ways `s[0:i]` contains `t[0:j]` as subsequence → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Rolling Array** — reduce space to O(n) by walking j from right to left

---

## Approaches Overview

| # | Approach              | Time   | Space  | When to use             |
|---|-----------------------|--------|--------|-------------------------|
| 1 | 2D Bottom-Up DP       | O(m·n) | O(m·n) | Clear; easy to reason   |
| 2 | 1D Rolling DP         | O(m·n) | O(n)   | Space-optimized version |

---

## Approach 1 — 2D Bottom-Up DP

### Intuition
`dp[i][j]` = number of distinct subsequences of `s[0:i]` that equal `t[0:j]`.

Two choices for `s[i-1]`:
- **Skip** it: `dp[i][j] += dp[i-1][j]`.
- **Use** it (only if `s[i-1] == t[j-1]`): `dp[i][j] += dp[i-1][j-1]`.

Base case: `dp[i][0] = 1` for all `i` — empty `t` is always matched (delete everything from `s`).

### Algorithm
1. `dp[i][0] = 1` for all i.
2. For i=1..m, j=1..n:
   - `dp[i][j] = dp[i-1][j]`.
   - If `s[i-1] == t[j-1]`: `dp[i][j] += dp[i-1][j-1]`.
3. Return `dp[m][n]`.

### Complexity
- **Time:** O(m·n)
- **Space:** O(m·n)

### Code
```go
func numDistinct(s string, t string) int {
    m, n := len(s), len(t)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
        dp[i][0] = 1
    }
    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            dp[i][j] = dp[i-1][j]
            if s[i-1] == t[j-1] { dp[i][j] += dp[i-1][j-1] }
        }
    }
    return dp[m][n]
}
```

### Dry Run
`s="rabbbit"`, `t="rabbit"` (abbreviated — 7×6 table)

Key cells:

| i\\j | "" | r | a | b | b | i | t |
|------|---|---|---|---|---|---|---|
| ""   | 1 | 0 | 0 | 0 | 0 | 0 | 0 |
| r    | 1 | 1 | 0 | 0 | 0 | 0 | 0 |
| a    | 1 | 1 | 1 | 0 | 0 | 0 | 0 |
| b    | 1 | 1 | 1 | 1 | 0 | 0 | 0 |
| b    | 1 | 1 | 1 | 2 | 1 | 0 | 0 |
| b    | 1 | 1 | 1 | 3 | 3 | 0 | 0 |
| i    | 1 | 1 | 1 | 3 | 3 | 3 | 0 |
| t    | 1 | 1 | 1 | 3 | 3 | 3 | 3 |

`dp[7][6] = 3` ✓

---

## Approach 2 — 1D Rolling DP

### Intuition
`dp[j]` depends only on the previous row's `dp[j]` and `dp[j-1]`. Roll into 1D, iterating `j` from right to left to avoid overwriting values we still need.

### Algorithm
1. `dp[0] = 1`, rest 0.
2. For each char in s:
   - For j from n down to 1:
     - If `s_char == t[j-1]`: `dp[j] += dp[j-1]`.

### Complexity
- **Time:** O(m·n)
- **Space:** O(n)

### Code
```go
func numDistinctOptimal(s string, t string) int {
    n := len(t)
    dp := make([]int, n+1)
    dp[0] = 1
    for _, sc := range s {
        for j := n; j >= 1; j-- {
            if sc == rune(t[j-1]) { dp[j] += dp[j-1] }
        }
    }
    return dp[n]
}
```

### Dry Run
`s="babgbag"`, `t="bag"`, dp=[1,0,0,0]:

| char | dp[3] | dp[2] | dp[1] |
|------|-------|-------|-------|
| b    | 0     | 0     | 1     |
| a    | 0     | 1     | 1     |
| b    | 0     | 1     | 2     |
| g    | 1     | 1     | 2     |
| b    | 1     | 1     | 3     |
| a    | 1     | 4     | 3     |
| g    | 5     | 4     | 3     |

`dp[3] = 5` ✓

---

## Key Takeaways
- Two-string DP: `dp[i][j]` depends on `dp[i-1][j]` (skip) and conditionally `dp[i-1][j-1]` (match).
- Right-to-left inner loop for 1D rolling ensures `dp[j-1]` used is from the previous outer iteration.
- `dp[i][0] = 1` is the crucial base case: empty pattern matched exactly once.

---

## Related Problems
- LeetCode #72 — Edit Distance
- LeetCode #1143 — Longest Common Subsequence
- LeetCode #583 — Delete Operation for Two Strings
