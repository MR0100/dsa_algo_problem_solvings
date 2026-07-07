# 0119 — Pascal's Triangle II

> LeetCode #119 · Difficulty: Easy
> **Categories:** Array, Dynamic Programming

---

## Problem Statement

Given an integer `rowIndex`, return the `rowIndex`-th (0-indexed) row of Pascal's triangle.

In Pascal's triangle, each number is the sum of the two numbers directly above it.

**Example 1:**
```
Input: rowIndex = 3
Output: [1,3,3,1]
```

**Example 2:**
```
Input: rowIndex = 0
Output: [1]
```

**Example 3:**
```
Input: rowIndex = 1
Output: [1,1]
```

**Constraints:**
- `0 <= rowIndex <= 33`

**Follow up:** Could you optimize your algorithm to use only O(rowIndex) extra space?

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Microsoft | ★★★☆☆ Medium | 2023          |
| Google    | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **In-place DP** — right-to-left update avoids overwriting needed values
- **Combinatorics** — C(n,k) = C(n,k-1) × (n-k+1) / k

---

## Approaches Overview

| # | Approach                | Time        | Space       | When to use           |
|---|-------------------------|-------------|-------------|-----------------------|
| 1 | In-Place Right-to-Left  | O(n²)       | O(n)        | General; satisfies follow-up |
| 2 | Combinatorial Formula   | O(n)        | O(n)        | Fastest                |

---

## Approach 1 — In-Place Right-to-Left Update

### Intuition
Maintain a single row of length `rowIndex+1`. Simulate building Pascal's triangle row by row, but update in-place by walking **right to left** so each `row[j] += row[j-1]` uses the previous row's `row[j-1]` (not yet overwritten).

### Algorithm
1. `row = [1, 0, 0, ..., 0]`.
2. For `i = 1..rowIndex`:
   - For `j = i..1` (right to left): `row[j] += row[j-1]`.

### Complexity
- **Time:** O(rowIndex²)
- **Space:** O(rowIndex)

### Code
```go
func getRow(rowIndex int) []int {
    row := make([]int, rowIndex+1)
    row[0] = 1
    for i := 1; i <= rowIndex; i++ {
        for j := i; j >= 1; j-- {
            row[j] += row[j-1]
        }
    }
    return row
}
```

### Dry Run
`rowIndex=3`:

| i | j walk     | row after         |
|---|-----------|-------------------|
| 1 | 1→1       | [1,1,0,0]         |
| 2 | 2→1       | [1,2,1,0]         |
| 3 | 3→1       | [1,3,3,1]         |

---

## Approach 2 — Combinatorial Formula

### Intuition
Row `n` entry `k` = C(n, k) = n! / (k! × (n-k)!). Use the recurrence:
`C(n, k) = C(n, k-1) × (n-k+1) / k`
to compute each entry from the previous in O(1).

### Complexity
- **Time:** O(rowIndex)
- **Space:** O(rowIndex)

### Code
```go
func getRowCombinatorial(rowIndex int) []int {
    row := make([]int, rowIndex+1)
    row[0] = 1
    for k := 1; k <= rowIndex; k++ {
        row[k] = row[k-1] * (rowIndex - k + 1) / k
    }
    return row
}
```

### Dry Run
`rowIndex=4`:

| k | row[k-1] × (n-k+1) / k | row[k] |
|---|------------------------|--------|
| 1 | 1 × 4 / 1              | 4      |
| 2 | 4 × 3 / 2              | 6      |
| 3 | 6 × 2 / 3              | 4      |
| 4 | 4 × 1 / 4              | 1      |

Result: [1, 4, 6, 4, 1] ✓

---

## Key Takeaways
- Right-to-left in-place update: same as the rolling-array trick used in 0-1 knapsack and distinct subsequences (#115).
- Combinatorial formula: O(n) time; note integer division must be exact at each step — multiply first before dividing.

---

## Related Problems
- LeetCode #118 — Pascal's Triangle (all rows)
