# 0054 — Spiral Matrix

> LeetCode #54 · Difficulty: Medium
> **Categories:** Array, Matrix, Simulation

---

## Problem Statement

Given an `m x n` `matrix`, return all elements of the `matrix` in **spiral order**.

**Example 1**
```
Input:  matrix = [[1,2,3],[4,5,6],[7,8,9]]
Output: [1,2,3,6,9,8,7,4,5]
```

**Example 2**
```
Input:  matrix = [[1,2,3,4],[5,6,7,8],[9,10,11,12]]
Output: [1,2,3,4,8,12,11,10,9,5,6,7]
```

**Constraints**
- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 10`
- `-100 <= matrix[i][j] <= 100`

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

- **Boundary Shrinking** — maintain four boundaries (top, bottom, left, right) and shrink each after traversing the corresponding edge.
- **Direction Vector Simulation** — walk in one direction; turn right when blocked or out of bounds.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Layer-by-Layer Peeling ✅ | O(m × n) | O(1) | Cleanest; preferred in interviews |
| 2 | Direction Vector Simulation | O(m × n) | O(m × n) | Handles any traversal pattern; visited array needed |

---

## Approach 1 — Layer-by-Layer Peeling (Recommended ✅)

### Intuition
Think of the matrix as concentric "rings." Each iteration, traverse the outermost ring (right along top, down right side, left along bottom, up left side) then shrink all four boundaries inward.

Guards are needed for the left and up traversals because after traversing right and down, the matrix may have collapsed to a single row or column.

### Algorithm
```
top=0, bottom=m-1, left=0, right=n-1
while top<=bottom and left<=right:
  right across top row [left..right];   top++
  down right col [top..bottom];         right--
  if top<=bottom: left across bottom [right..left]; bottom--
  if left<=right: up left col [bottom..top];        left++
```

### Complexity
- **Time:** O(m × n) — every element visited exactly once.
- **Space:** O(1) — only four boundary pointers (output list not counted).

### Code
```go
func layerPeel(matrix [][]int) []int {
    m, n := len(matrix), len(matrix[0])
    result := []int{}
    top, bottom, left, right := 0, m-1, 0, n-1
    for top <= bottom && left <= right {
        for c := left; c <= right; c++  { result = append(result, matrix[top][c]) }; top++
        for r := top; r <= bottom; r++  { result = append(result, matrix[r][right]) }; right--
        if top <= bottom {
            for c := right; c >= left; c-- { result = append(result, matrix[bottom][c]) }; bottom--
        }
        if left <= right {
            for r := bottom; r >= top; r-- { result = append(result, matrix[r][left]) }; left++
        }
    }
    return result
}
```

### Dry Run — `matrix = [[1,2,3],[4,5,6],[7,8,9]]`
```
Initial: top=0, bottom=2, left=0, right=2

Iteration 1:
  Right across top (row 0): [1, 2, 3]; top=1
  Down right (col 2, rows 1–2): [6, 9]; right=1
  Left across bottom (row 2, cols 1–0): [8, 7]; bottom=1
  Up left (col 0, rows 1–1): [4]; left=1

Result so far: [1,2,3,6,9,8,7,4]

Iteration 2: top=1, bottom=1, left=1, right=1
  Right across top (row 1, col 1): [5]; top=2
  Down right (col 1, rows 2..1): nothing (top>bottom)
  top>bottom → skip bottom traversal
  left≤right but bottom<top → skip left traversal

Result: [1,2,3,6,9,8,7,4,5] ✓
```

---

## Approach 2 — Direction Vector Simulation

### Intuition
Walk one step at a time in the current direction (right, down, left, up). Mark each cell visited. When the next cell is out of bounds or already visited, rotate 90° clockwise. Record as we walk.

### Complexity
- **Time:** O(m × n).
- **Space:** O(m × n) — visited array.

---

## Key Takeaways

- **Guards for inner traversals** — after traversing right and down, the remaining `top` might exceed `bottom` (single-row matrix or collapsed). Always guard `if top <= bottom` before the bottom-left traversal, and `if left <= right` before the up traversal.
- **Direction simulation is more general** — can handle any traversal pattern (diagonal, zigzag) by changing the direction array. Layer peeling only works for strict spiral.
- **Companion problem** — #59 Spiral Matrix II fills the matrix in spiral order; the same layer-peel framework works but writes indices instead of reading values.

---

## Related Problems

- LeetCode #59 — Spiral Matrix II (fill matrix in spiral order)
- LeetCode #48 — Rotate Image (matrix manipulation)
- LeetCode #73 — Set Matrix Zeroes (in-place matrix update)
