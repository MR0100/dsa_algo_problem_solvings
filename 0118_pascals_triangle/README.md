# 0118 — Pascal's Triangle

> LeetCode #118 · Difficulty: Easy
> **Categories:** Array, Dynamic Programming

---

## Problem Statement

Given an integer `numRows`, return the first `numRows` of Pascal's triangle.

In Pascal's triangle, each number is the sum of the two numbers directly above it.

**Example 1:**
```
Input: numRows = 5
Output: [[1],[1,1],[1,2,1],[1,3,3,1],[1,4,6,4,1]]
```

**Example 2:**
```
Input: numRows = 1
Output: [[1]]
```

**Constraints:**
- `1 <= numRows <= 30`

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Apple     | ★★★☆☆ Medium | 2023          |
| Google    | ★★★☆☆ Medium | 2023          |
| Microsoft | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **DP / Recurrence** — `triangle[i][j] = triangle[i-1][j-1] + triangle[i-1][j]`

---

## Approaches Overview

| # | Approach     | Time         | Space        | When to use |
|---|--------------|--------------|--------------|-------------|
| 1 | Iterative    | O(numRows²)  | O(numRows²)  | Always      |

---

## Approach 1 — Iterative Row-by-Row

### Intuition
Each row starts and ends with 1. Interior element `[i][j] = prev[j-1] + prev[j]`. Build each row from the previous.

### Algorithm
1. For `i = 0..numRows-1`:
   - Allocate row of length `i+1`.
   - `row[0] = row[i] = 1`.
   - For `j = 1..i-1`: `row[j] = result[i-1][j-1] + result[i-1][j]`.
2. Append row to result.

### Complexity
- **Time:** O(numRows²) — 1+2+...+numRows = numRows(numRows+1)/2 elements.
- **Space:** O(numRows²) — all elements stored in result.

### Code
```go
func generate(numRows int) [][]int {
    result := make([][]int, numRows)
    for i := 0; i < numRows; i++ {
        row := make([]int, i+1)
        row[0] = 1; row[i] = 1
        for j := 1; j < i; j++ {
            row[j] = result[i-1][j-1] + result[i-1][j]
        }
        result[i] = row
    }
    return result
}
```

### Dry Run
`numRows=5`:

| row | values         |
|-----|----------------|
| 0   | [1]            |
| 1   | [1, 1]         |
| 2   | [1, 1+1=2, 1]  |
| 3   | [1, 1+2=3, 2+1=3, 1] |
| 4   | [1, 1+3=4, 3+3=6, 3+1=4, 1] |

---

## Key Takeaways
- First and last elements of every row are always 1.
- Interior: `row[j] = prevRow[j-1] + prevRow[j]`.
- To return only one row: see #119 (O(k) space via in-place update right-to-left).

---

## Related Problems
- LeetCode #119 — Pascal's Triangle II (return specific row)
