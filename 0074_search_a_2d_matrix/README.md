# 0074 — Search a 2D Matrix

> LeetCode #74 · Difficulty: Medium
> **Categories:** Array, Binary Search, Matrix

---

## Problem Statement

You are given an `m x n` integer matrix `matrix` with the following two properties:
- Each row is sorted in non-decreasing order.
- The first integer of each row is greater than the last integer of the previous row.

Given an integer `target`, return `true` if `target` is in `matrix` or `false` otherwise.

You must write a solution in `O(log(m * n))` time complexity.

**Example 1**
```
Input:  matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 3
Output: true
```

**Example 2**
```
Input:  matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]], target = 13
Output: false
```

**Constraints**
- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 100`
- `-10⁴ <= matrix[i][j], target <= 10⁴`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search** — treat 2D matrix as a flat sorted array; index mapping `(mid) → (mid/n, mid%n)`.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two-Step Binary Search | O(log m + log n) | O(1) | Find row first, then search row |
| 2 | Flat Binary Search ✅ | O(log(m × n)) | O(1) | Cleaner; treats matrix as 1D sorted array |

---

## Approach 1 — Two-Step Binary Search

### Intuition
1. Binary search rows to find the last row where `matrix[row][0] ≤ target`.
2. Binary search within that row.

### Complexity
- **Time:** O(log m + log n) = O(log(m × n)).
- **Space:** O(1).

---

## Approach 2 — Flat Binary Search (Recommended ✅)

### Intuition
The matrix property guarantees that if we read row by row, the elements form a strictly increasing sequence. So the entire matrix is one sorted array of `m × n` elements. Binary search on virtual indices `[0, m×n-1]`:
- `matrix[mid / n][mid % n]` converts flat index to 2D.

### Algorithm
```
lo=0, hi=m*n-1
while lo<=hi:
  mid=(lo+hi)/2; val=matrix[mid/n][mid%n]
  if val==target: true
  elif val<target: lo=mid+1
  else: hi=mid-1
false
```

### Complexity
- **Time:** O(log(m × n)).
- **Space:** O(1).

### Code
```go
func flatBinarySearch(matrix [][]int, target int) bool {
    m, n := len(matrix), len(matrix[0])
    lo, hi := 0, m*n-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        val := matrix[mid/n][mid%n]
        if val == target { return true }
        if val < target { lo = mid+1 } else { hi = mid-1 }
    }
    return false
}
```

### Dry Run — `matrix = [[1,3,5,7],[10,11,16,20],[23,30,34,60]]`, `target = 3`
```
m=3, n=4, m*n=12. lo=0, hi=11.

mid=5: matrix[5/4][5%4] = matrix[1][1] = 11 > 3 → hi=4.
mid=2: matrix[2/4][2%4] = matrix[0][2] = 5 > 3 → hi=1.
mid=0: matrix[0][0] = 1 < 3 → lo=1.
mid=1: matrix[0][1] = 3 == 3 → return true ✓
```

---

## Key Takeaways

- **`mid / n` and `mid % n`** — the core insight. Integer division gives the row; modulo gives the column.
- **This only works for Approach 2 because of the strict sorted-row property** — if the matrix were sorted within rows but NOT across rows, this flat-array trick would fail. For that case (#240), use the staircase search instead.
- **Compare with #240 (Search a 2D Matrix II)** — #74 has the global sorted property → binary search. #240 does not → staircase search from top-right corner.

---

## Related Problems

- LeetCode #240 — Search a 2D Matrix II (sorted rows and columns; staircase search)
- LeetCode #35 — Search Insert Position (lower-bound binary search)
- LeetCode #162 — Find Peak Element (binary search with non-sorted property)
