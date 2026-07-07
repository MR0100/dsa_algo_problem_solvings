# 0304 — Range Sum Query 2D - Immutable

> LeetCode #304 · Difficulty: Medium
> **Categories:** Array, Matrix, Design, Prefix Sum

---

## Problem Statement

Given a 2D matrix `matrix`, handle multiple queries of the following type:

- Calculate the **sum** of the elements of `matrix` inside the rectangle defined by its **upper left corner** `(row1, col1)` and **lower right corner** `(row2, col2)`.

Implement the `NumMatrix` class:

- `NumMatrix(int[][] matrix)` Initializes the object with the integer matrix `matrix`.
- `int sumRegion(int row1, int col1, int row2, int col2)` Returns the **sum** of the elements of `matrix` inside the rectangle defined by its **upper left corner** `(row1, col1)` and **lower right corner** `(row2, col2)`.

You must design an algorithm where `sumRegion` works on `O(1)` time complexity.

**Example 1:**

```
Input
["NumMatrix", "sumRegion", "sumRegion", "sumRegion"]
[[[[3, 0, 1, 4, 2], [5, 6, 3, 2, 1], [1, 2, 0, 1, 5], [4, 1, 0, 1, 7], [1, 0, 3, 0, 5]]], [2, 1, 4, 3], [1, 1, 2, 2], [1, 2, 2, 4]]
Output
[null, 8, 11, 12]

Explanation
NumMatrix numMatrix = new NumMatrix([[3, 0, 1, 4, 2], [5, 6, 3, 2, 1], [1, 2, 0, 1, 5], [4, 1, 0, 1, 7], [1, 0, 3, 0, 5]]);
numMatrix.sumRegion(2, 1, 4, 3); // return 8 (i.e sum of the red rectangle)
numMatrix.sumRegion(1, 1, 2, 2); // return 11 (i.e sum of the green rectangle)
numMatrix.sumRegion(1, 2, 2, 4); // return 12 (i.e sum of the blue rectangle)
```

**Constraints:**

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 200`
- `-10^4 <= matrix[i][j] <= 10^4`
- `0 <= row1 <= row2 < m`
- `0 <= col1 <= col2 < n`
- At most `10^4` calls will be made to `sumRegion`.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★★☆ High       | 2024          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Facebook  | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **2D Prefix Sum (integral image)** — cumulative area sums answer any sub-rectangle in O(1) via inclusion–exclusion → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Matrix traversal** — building the cumulative grid in row-major order → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Design (immutable structure)** — heavy one-time preprocessing for cheap repeated queries → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | Constructor | Query | Space | When to use |
|---|----------|-------------|-------|-------|-------------|
| 1 | Brute force (re-sum) | O(1) | O(m·n) | O(m·n) | Few queries or correctness baseline |
| 2 | 2D prefix sum (Optimal) | O(m·n) | O(1) | O(m·n) | Many queries on an immutable matrix — required for O(1) queries |

---

## Approach 1 — Brute Force (Re-sum Each Query)

### Intuition
Store the matrix; for each query, double-loop over the requested rows and columns adding every cell. Correct but O(rows·cols) per query, which fails the O(1)-query requirement.

### Algorithm
1. Keep the matrix reference (immutable, so no copy).
2. For `sumRegion`, loop `r` from `row1` to `row2` and `c` from `col1` to `col2`, accumulating `matrix[r][c]`.
3. Return the sum.

### Complexity
- **Time:** constructor O(1); each query O((row2−row1+1)·(col2−col1+1)).
- **Space:** O(m·n) to hold the matrix.

### Code
```go
type NumMatrixBrute struct {
	matrix [][]int // the original, immutable grid
}

func NewNumMatrixBrute(matrix [][]int) NumMatrixBrute {
	return NumMatrixBrute{matrix: matrix}
}

func (m NumMatrixBrute) SumRegion(row1, col1, row2, col2 int) int {
	sum := 0
	for r := row1; r <= row2; r++ { // each row of the rectangle
		for c := col1; c <= col2; c++ { // each column of the rectangle
			sum += m.matrix[r][c]
		}
	}
	return sum
}
```

### Dry Run
Query `sumRegion(2, 1, 4, 3)` on the example matrix.

| Row r | cells matrix[r][1..3] | row sum | running |
|---|---|---|---|
| 2 | 2, 0, 1 | 3 | 3 |
| 3 | 1, 0, 1 | 2 | 5 |
| 4 | 0, 3, 0 | 3 | 8 |

Returns **8**.

---

## Approach 2 — 2D Prefix Sum (Optimal)

### Intuition
Let `P[r][c]` be the sum of the whole sub-rectangle from `(0,0)` to `(r−1, c−1)` (an "integral image"). Any rectangle sum comes from four corners by inclusion–exclusion:

```
sum(r1,c1,r2,c2) = P[r2+1][c2+1] − P[r1][c2+1] − P[r2+1][c1] + P[r1][c1]
```

We subtract the strip above and the strip to the left, then add back the top-left corner that got subtracted twice. A padded all-zero first row/column removes boundary cases.

### Algorithm
1. Allocate `P` of size `(m+1)×(n+1)`, all zeros (the padding row/column).
2. Fill `P[r+1][c+1] = matrix[r][c] + P[r][c+1] + P[r+1][c] − P[r][c]`.
3. Answer each query with the four-corner formula above.

### Complexity
- **Time:** constructor O(m·n); each `SumRegion` O(1).
- **Space:** O(m·n) for the prefix grid.

### Code
```go
type NumMatrix struct {
	prefix [][]int // (m+1)×(n+1) padded cumulative-area sums
}

func NewNumMatrix(matrix [][]int) NumMatrix {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return NumMatrix{prefix: [][]int{{0}}}
	}
	m, n := len(matrix), len(matrix[0])
	prefix := make([][]int, m+1) // extra top row / left column of zeros
	for i := range prefix {
		prefix[i] = make([]int, n+1)
	}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			prefix[r+1][c+1] = matrix[r][c] +
				prefix[r][c+1] +
				prefix[r+1][c] -
				prefix[r][c]
		}
	}
	return NumMatrix{prefix: prefix}
}

func (nm NumMatrix) SumRegion(row1, col1, row2, col2 int) int {
	p := nm.prefix
	return p[row2+1][col2+1] - // full area to bottom-right corner
		p[row1][col2+1] - // remove the strip above the rectangle
		p[row2+1][col1] + // remove the strip to the left
		p[row1][col1] //     add back the top-left, subtracted twice
}
```

### Dry Run
Query `sumRegion(2, 1, 4, 3)`.

First a few relevant prefix values (P is 1-indexed over the padded grid). The full area to bottom-right corner `P[5][4]`, the strip above `P[2][4]`, the strip left `P[5][1]`, and the corner `P[2][1]` combine as:

```
sum = P[5][4] − P[2][4] − P[5][1] + P[2][1]
```

Computing cumulative sums over the example matrix yields `P[5][4] = 37`, `P[2][4] = 27`, `P[5][1] = 14`, `P[2][1] = 8`:

```
sum = 37 − 27 − 14 + 8 = 8
```

Returns **8** — matching the brute-force trace above.

---

## Key Takeaways

- The **2D prefix sum / integral image** is the canonical trick for O(1) rectangle queries; it generalizes the 1D prefix sum by inclusion–exclusion over four corners.
- Build recurrence: `P[r+1][c+1] = cell + above + left − overlap`. Query: `P[br] − P[top] − P[left] + P[corner]`.
- The **+1 zero padding** eliminates all boundary checks — no special-casing `row1 = 0` or `col1 = 0`.
- Preprocessing is O(m·n) and pays for itself the moment you have more than a couple of queries.

---

## Related Problems

- LeetCode #303 — Range Sum Query - Immutable (1D version)
- LeetCode #308 — Range Sum Query 2D - Mutable (2D Fenwick tree)
- LeetCode #221 — Maximal Square (2D DP on a grid)
- LeetCode #1074 — Number of Submatrices That Sum to Target (2D prefix + hashing)
