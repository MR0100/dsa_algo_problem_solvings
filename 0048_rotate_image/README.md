# 0048 — Rotate Image

> LeetCode #48 · Difficulty: Medium
> **Categories:** Array, Math, Matrix

---

## Problem Statement

You are given an `n x n` 2D `matrix` representing an image, rotate the image by **90 degrees** (clockwise).

You have to rotate the image **in-place**, which means you have to modify the input 2D matrix directly. **Do not** allocate another 2D matrix and do the rotation.

**Example 1**
```
Input:  matrix = [[1,2,3],[4,5,6],[7,8,9]]
Output:           [[7,4,1],[8,5,2],[9,6,3]]
```

**Example 2**
```
Input:  matrix = [[5,1,9,11],[2,4,8,10],[13,3,6,7],[15,14,12,16]]
Output:           [[15,13,2,5],[14,3,4,1],[12,6,8,9],[16,7,10,11]]
```

**Constraints**
- `n == matrix.length == matrix[i].length`
- `1 <= n <= 20`
- `-1000 <= matrix[i][j] <= 1000`

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

- **Matrix Transformation** — transpose (reflect across main diagonal) + row reversal = 90° clockwise rotation.
- **In-Place Swapping** — only two pointer variables needed; no extra matrix.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Extra Matrix | O(n²) | O(n²) | Violates in-place requirement; reference only |
| 2 | Transpose + Reverse Rows ✅ | O(n²) | O(1) | The standard in-place answer |

---

## Approach 1 — Extra Matrix

### Intuition
For a 90° clockwise rotation: element at `(r, c)` moves to `(c, n-1-r)`. Use an extra matrix.

### Complexity
- **Time:** O(n²).
- **Space:** O(n²).

---

## Approach 2 — Transpose + Reverse Rows (Recommended ✅)

### Intuition
A 90° clockwise rotation decomposes into two reflections:

1. **Transpose** (reflect across main diagonal): `(r, c) ↔ (c, r)`.
2. **Reverse each row** (reflect across the vertical midline).

Combined:
```
Original:   Transpose:   Reverse rows:
1 2 3       1 4 7        7 4 1
4 5 6  →    2 5 8    →   8 5 2
7 8 9       3 6 9        9 6 3
```

### Why it works
After transposing, column `j` of the original becomes row `j`. Reversing each row then puts the elements in the correct order for a clockwise rotation.

**Counter-clockwise = transpose + reverse columns.**

### Algorithm
```
Transpose: for r=0 to n-1, c=r+1 to n-1: swap(matrix[r][c], matrix[c][r])
Reverse rows: for r=0 to n-1: reverse matrix[r]
```

### Complexity
- **Time:** O(n²) — n²/2 swaps for transpose + n²/2 swaps for row reversal.
- **Space:** O(1) — in-place.

### Code
```go
func rotateInPlace(matrix [][]int) {
    n := len(matrix)
    for r := 0; r < n; r++ {
        for c := r+1; c < n; c++ { matrix[r][c], matrix[c][r] = matrix[c][r], matrix[r][c] }
    }
    for r := 0; r < n; r++ {
        left, right := 0, n-1
        for left < right { matrix[r][left], matrix[r][right] = matrix[r][right], matrix[r][left]; left++; right-- }
    }
}
```

### Dry Run — `matrix = [[1,2,3],[4,5,6],[7,8,9]]`
```
Transpose:
  swap(0,1),(0,2): [1,4,7]
  swap(0,1),(0,2): [1,4,7],[2,5,8],[3,6,9]

After transpose:
  1 4 7
  2 5 8
  3 6 9

Reverse each row:
  Row 0: [7,4,1]
  Row 1: [8,5,2]
  Row 2: [9,6,3]

Result: [[7,4,1],[8,5,2],[9,6,3]] ✓
```

---

## Key Takeaways

- **Two reflections = one rotation** — any rotation can be decomposed into two reflections. 90° CW = transpose + horizontal flip; 90° CCW = transpose + vertical flip; 180° = horizontal flip + vertical flip (or just rotate each element).
- **Transpose only swaps `c > r`** — start the inner loop at `c = r+1` to avoid swapping elements back.
- **4-element cycle alternative** — another O(1)-space approach rotates 4 elements in a cycle at once, layer by layer. Correct but harder to implement without mistakes; transpose+reverse is simpler.

---

## Related Problems

- LeetCode #54 — Spiral Matrix (another matrix traversal pattern)
- LeetCode #73 — Set Matrix Zeroes (in-place matrix modification)
- LeetCode #289 — Game of Life (in-place matrix update with encoding trick)
