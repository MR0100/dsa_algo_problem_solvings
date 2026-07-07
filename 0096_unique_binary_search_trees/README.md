# 0096 — Unique Binary Search Trees

> LeetCode #96 · Difficulty: Medium
> **Categories:** Math, Dynamic Programming, Tree, Binary Search Tree

---

## Problem Statement

Given an integer `n`, return the number of structurally unique **BST's** (binary search trees) which has exactly `n` nodes of unique values from `1` to `n`.

**Example 1:**
```
Input: n = 3
Output: 5
```

**Example 2:**
```
Input: n = 1
Output: 1
```

**Constraints:**
- `1 <= n <= 19`

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★☆☆ Medium   | 2024          |
| Google    | ★★★☆☆ Medium   | 2023          |
| Microsoft | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming** — `dp[n]` = Catalan number C(n). See [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Catalan Numbers** — C(n) = Σ C(i-1)×C(n-i) for i=1..n; closed form C(n) = C(2n,n)/(n+1).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP (bottom-up) | O(n²) | O(n) | Standard; avoids overflow risks |
| 2 | Catalan Formula | O(n) | O(1) | Elegant closed-form |

---

## Approach 1 — DP (Catalan Number via Recurrence)

### Intuition
For `n` nodes, any node `i` (1 ≤ i ≤ n) can be the root. The left subtree contains `i-1` nodes and the right subtree contains `n-i` nodes. The number of unique BSTs with `i` as root = `dp[i-1] × dp[n-i]`. Summing over all roots gives `dp[n]`.

**Catalan recurrence:** `dp[n] = Σ(j=1..n) dp[j-1] × dp[n-j]`.

### Algorithm
1. `dp[0] = 1, dp[1] = 1`.
2. For `i = 2` to `n`: `dp[i] = Σ(j=1..i) dp[j-1] × dp[i-j]`.

### Complexity
- **Time:** O(n²) — n values, each computed by O(n) sum.
- **Space:** O(n)

### Code
```go
func numTrees(n int) int {
    dp := make([]int, n+1)
    dp[0], dp[1] = 1, 1
    for i := 2; i <= n; i++ {
        for j := 1; j <= i; j++ {
            dp[i] += dp[j-1] * dp[i-j]
        }
    }
    return dp[n]
}
```

### Dry Run (n=3)

| i | j | dp[j-1] × dp[i-j] | dp[i] |
|---|---|-------------------|-------|
| 2 | 1 | dp[0]×dp[1]=1 | 1 |
| 2 | 2 | dp[1]×dp[0]=1 | 2 |
| 3 | 1 | dp[0]×dp[2]=2 | 2 |
| 3 | 2 | dp[1]×dp[1]=1 | 3 |
| 3 | 3 | dp[2]×dp[0]=2 | 5 |

dp[3] = 5 ✓

---

## Approach 2 — Catalan Number Formula

### Intuition
The n-th Catalan number has a closed form: `C(n) = C(2n,n) / (n+1)`.

Computed iteratively to avoid overflow: at step `i` (0..n-1), multiply by `(n+1+i)` and divide by `(i+1)`. Division is exact at each step (the intermediate result is always an integer).

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func numTreesFormula(n int) int {
    result := 1
    for i := 0; i < n; i++ {
        result = result * (n + 1 + i) / (i + 1)
    }
    return result / (n + 1)
}
```

### Dry Run (n=3)

Init: `result = 1`. Loop `i = 0..2`, updating `result = result * (n+1+i) / (i+1)` (n=3):

| i | n+1+i | i+1 | result * (n+1+i) | / (i+1) → result |
|---|-------|-----|------------------|------------------|
| 0 | 4     | 1   | 1 × 4 = 4        | 4 / 1 = 4        |
| 1 | 5     | 2   | 4 × 5 = 20       | 20 / 2 = 10      |
| 2 | 6     | 3   | 10 × 6 = 60      | 60 / 3 = 20      |

After loop: `result = 20`. Final: `result / (n+1) = 20 / 4 = 5` ✓

### Catalan Sequence (first 10)

| n | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 |
|---|---|---|---|---|---|---|---|---|---|---|
| C(n) | 1 | 1 | 2 | 5 | 14 | 42 | 132 | 429 | 1430 | 4862 |

---

## Key Takeaways
- `dp[0] = 1` represents the empty tree — critical base case for the product formula.
- Catalan numbers appear in: BST counting, valid parentheses counting, triangulations, mountain ranges.
- The DP recurrence is O(n²) but since n ≤ 19 (constraint), this is fine.
- Memorise Catalan(n): 1,1,2,5,14,42,132... Many "how many structures" problems give Catalan answers.

---

## Related Problems
- LeetCode #95 — Unique Binary Search Trees II (enumerate all trees, not just count)
- LeetCode #22 — Generate Parentheses (count = Catalan number)
- LeetCode #241 — Different Ways to Add Parentheses (same recurrence structure)
