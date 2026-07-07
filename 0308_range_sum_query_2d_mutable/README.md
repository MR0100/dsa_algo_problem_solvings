# 0308 — Range Sum Query 2D - Mutable

> LeetCode #308 · Difficulty: Medium (Premium)
> **Categories:** Array, Matrix, Design, Binary Indexed Tree, Segment Tree

---

## Problem Statement

Given a 2D matrix `matrix`, handle multiple queries of the following types:

1. **Update** the value of a cell in `matrix`.
2. Calculate the **sum** of the elements of `matrix` inside the rectangle defined by its **upper left corner** `(row1, col1)` and **lower right corner** `(row2, col2)`.

Implement the `NumMatrix` class:

- `NumMatrix(int[][] matrix)` Initializes the object with the integer matrix `matrix`.
- `void update(int row, int col, int val)` **Updates** the value of `matrix[row][col]` to be `val`.
- `int sumRegion(int row1, int col1, int row2, int col2)` Returns the **sum** of the elements of `matrix` inside the rectangle defined by its **upper left corner** `(row1, col1)` and **lower right corner** `(row2, col2)`.

**Example 1:**

```
Input:
["NumMatrix", "sumRegion", "update", "sumRegion"]
[[[[3, 0, 1, 4, 2], [5, 6, 3, 2, 1], [1, 2, 0, 1, 5], [4, 1, 0, 1, 7], [1, 0, 3, 0, 5]]], [2, 1, 4, 3], [3, 2, 2], [2, 1, 4, 3]]
Output:
[null, 8, null, 10]

Explanation:
NumMatrix numMatrix = new NumMatrix([[3, 0, 1, 4, 2], [5, 6, 3, 2, 1], [1, 2, 0, 1, 5], [4, 1, 0, 1, 7], [1, 0, 3, 0, 5]]);
numMatrix.sumRegion(2, 1, 4, 3); // return 8
numMatrix.update(3, 2, 2);       // matrix[3][2] becomes 2
numMatrix.sumRegion(2, 1, 4, 3); // return 10
```

**Constraints:**

- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 200`
- `-10^5 <= matrix[i][j] <= 10^5`
- `0 <= row < m`
- `0 <= col < n`
- `-10^5 <= val <= 10^5`
- `0 <= row1 <= row2 < m`
- `0 <= col1 <= col2 < n`
- At most `5000` calls will be made to `sumRegion` and `update`.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2024          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★★☆☆ Medium     | 2023          |
| Apple     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **2D Binary Indexed Tree (Fenwick)** — nested lowbit ascent/descent for O(log m · log n) update and prefix-rectangle sum → see [`/dsa/segment_tree_fenwick.md`](/dsa/segment_tree_fenwick.md)
- **2D Prefix Sums** — per-row prefixes and inclusion–exclusion of four corner prefixes → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Matrix Traversal** — iterating rectangles and rows over a grid → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Design Data Structures** — a class with a fixed update/query contract → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | Update | SumRegion | Space | When to use |
|---|----------|--------|-----------|-------|-------------|
| 1 | Brute Force (matrix) | O(1) | O(m·n) | O(m·n) | Very few queries |
| 2 | Per-Row Prefix Sums | O(n) | O(rows) | O(m·n) | Wide-short grids; balanced ops |
| 3 | 2D Fenwick Tree (Optimal) | O(log m·log n) | O(log m·log n) | O(m·n) | Interleaved update + region sum |

---

## Approach 1 — Brute Force

### Intuition

Store the grid. `update` is one assignment; `sumRegion` loops over every cell in the rectangle. Correctness baseline that exposes the cost the smarter structures remove.

### Algorithm

1. Deep-copy the matrix.
2. `Update(row, col, val)`: `mat[row][col] = val`.
3. `SumRegion(...)`: double loop over `row1..row2` × `col1..col2` adding each cell.

### Complexity

- **Time:** Update O(1); SumRegion O(m·n) (rectangle area).
- **Space:** O(m·n) — the stored matrix.

### Code

```go
type BruteForce struct {
	mat [][]int
}

func NewBruteForce(matrix [][]int) *BruteForce {
	m := len(matrix)
	cp := make([][]int, m)
	for i := range matrix {
		cp[i] = append([]int(nil), matrix[i]...)
	}
	return &BruteForce{mat: cp}
}

func (b *BruteForce) Update(row, col, val int) {
	b.mat[row][col] = val
}

func (b *BruteForce) SumRegion(row1, col1, row2, col2 int) int {
	sum := 0
	for r := row1; r <= row2; r++ {
		for c := col1; c <= col2; c++ {
			sum += b.mat[r][c]
		}
	}
	return sum
}
```

### Dry Run

Region `(2,1)..(4,3)` on the initial grid — rows 2,3,4 and cols 1,2,3:

| Row | Cells (col 1..3) | Row sum |
|-----|------------------|---------|
| 2 | 2, 0, 1 | 3 |
| 3 | 1, 0, 1 | 2 |
| 4 | 0, 3, 0 | 3 |

Total = 3 + 2 + 3 = **8**. After `update(3,2,2)` row 3 becomes `1, 2, 1` → row sum 4; total = 3 + 4 + 3 = **10**.

---

## Approach 2 — Per-Row Prefix Sums

### Intuition

Balance update against query. Each row keeps prefix sums, so a single row segment `[col1, col2]` is O(1). A region loops the rows `row1..row2` adding each row's segment. An update rewrites the suffix of one row's prefixes. Good for wide-but-short grids.

### Algorithm

1. Build `rowPre[r][c+1] = rowPre[r][c] + matrix[r][c]`.
2. `Update(row, col, val)`: set the value; recompute `rowPre[row][col+1..n]`.
3. `SumRegion(...)`: for each row `r` in the region, add `rowPre[r][col2+1] − rowPre[r][col1]`.

### Complexity

- **Time:** Update O(n); SumRegion O(rows in region).
- **Space:** O(m·n) — a prefix array per row plus the matrix.

### Code

```go
type RowPrefix struct {
	mat    [][]int
	rowPre [][]int
	cols   int
}

func NewRowPrefix(matrix [][]int) *RowPrefix {
	m := len(matrix)
	n := 0
	if m > 0 {
		n = len(matrix[0])
	}
	mat := make([][]int, m)
	rowPre := make([][]int, m)
	for r := 0; r < m; r++ {
		mat[r] = append([]int(nil), matrix[r]...)
		rowPre[r] = make([]int, n+1)
		for c := 0; c < n; c++ {
			rowPre[r][c+1] = rowPre[r][c] + matrix[r][c]
		}
	}
	return &RowPrefix{mat: mat, rowPre: rowPre, cols: n}
}

func (p *RowPrefix) Update(row, col, val int) {
	p.mat[row][col] = val
	for c := col; c < p.cols; c++ {
		p.rowPre[row][c+1] = p.rowPre[row][c] + p.mat[row][c]
	}
}

func (p *RowPrefix) SumRegion(row1, col1, row2, col2 int) int {
	sum := 0
	for r := row1; r <= row2; r++ {
		sum += p.rowPre[r][col2+1] - p.rowPre[r][col1]
	}
	return sum
}
```

### Dry Run

Row prefixes for rows 2..4 (length 6, index 0 = 0):

| Row | rowPre |
|-----|--------|
| 2 (`1 2 0 1 5`) | 0,1,3,3,4,9 |
| 3 (`4 1 0 1 7`) | 0,4,5,5,6,13 |
| 4 (`1 0 3 0 5`) | 0,1,1,4,4,9 |

`sumRegion(2,1,4,3)`: row2 seg = `rowPre[2][4]−rowPre[2][1]` = 4−1 = 3; row3 = 6−4 = 2; row4 = 4−1 = 3; total **8**. After `update(3,2,2)` row 3 prefixes → 0,4,5,7,8,15; row3 seg = 8−4 = 4; total = 3+4+3 = **10**.

---

## Approach 3 — 2D Fenwick Tree (Optimal)

### Intuition

Extend the 1D Fenwick tree to two dimensions. `tree[i][j]` holds a partial sum over a rectangle sized by the lowbits of `i` and `j`. A prefix sum over `(0,0)..(r,c)` is a nested lowbit descent; a point update is a nested lowbit ascent. A general region uses inclusion–exclusion:

```
sum(r1,c1,r2,c2) = P(r2,c2) − P(r1−1,c2) − P(r2,c1−1) + P(r1−1,c1−1)
```

Because updates set (not add), cache the grid and push the delta `val − old`.

### Algorithm

1. Build all zeros; `add` each initial value via the delta path.
2. `add(r,c,d)`: `for i=r..m by lowbit: for j=c..n by lowbit: tree[i][j] += d`.
3. `query(r,c)`: `for i=r..1 by lowbit: for j=c..1 by lowbit: s += tree[i][j]`.
4. `Update(row,col,val)`: `d = val − nums[row][col]`; store `val`; `add(row+1,col+1,d)`.
5. `SumRegion(...)`: inclusion–exclusion of four `query` calls (1-indexed corners).

### Complexity

- **Time:** Update O(log m · log n), SumRegion O(log m · log n); Build O(mn log m log n).
- **Space:** O(m·n) — the tree plus the value cache.

### Code

```go
type FenwickTree2D struct {
	m, n int
	tree [][]int
	nums [][]int
}

func NewFenwickTree2D(matrix [][]int) *FenwickTree2D {
	m := len(matrix)
	n := 0
	if m > 0 {
		n = len(matrix[0])
	}
	f := &FenwickTree2D{m: m, n: n}
	f.tree = make([][]int, m+1)
	for i := range f.tree {
		f.tree[i] = make([]int, n+1)
	}
	f.nums = make([][]int, m)
	for r := 0; r < m; r++ {
		f.nums[r] = make([]int, n)
		for c := 0; c < n; c++ {
			f.nums[r][c] = 0
			f.Update(r, c, matrix[r][c])
		}
	}
	return f
}

func (f *FenwickTree2D) add(r, c, delta int) {
	for i := r; i <= f.m; i += i & (-i) {
		for j := c; j <= f.n; j += j & (-j) {
			f.tree[i][j] += delta
		}
	}
}

func (f *FenwickTree2D) query(r, c int) int {
	s := 0
	for i := r; i > 0; i -= i & (-i) {
		for j := c; j > 0; j -= j & (-j) {
			s += f.tree[i][j]
		}
	}
	return s
}

func (f *FenwickTree2D) Update(row, col, val int) {
	delta := val - f.nums[row][col]
	f.nums[row][col] = val
	f.add(row+1, col+1, delta)
}

func (f *FenwickTree2D) SumRegion(row1, col1, row2, col2 int) int {
	return f.query(row2+1, col2+1) -
		f.query(row1, col2+1) -
		f.query(row2+1, col1) +
		f.query(row1, col1)
}
```

### Dry Run

Corners of `sumRegion(2,1,4,3)` via inclusion–exclusion (1-indexed):

| Term | query args | meaning |
|------|-----------|---------|
| +P(5,4) | prefix over rows 0..4, cols 0..3 | big rectangle |
| −P(2,4) | rows 0..1, cols 0..3 | strip above |
| −P(5,1) | rows 0..4, col 0 | strip left |
| +P(2,1) | rows 0..1, col 0 | double-subtracted corner added back |

Evaluating on the initial grid gives `P(5,4) − P(2,4) − P(5,1) + P(2,1) = 8`. After `update(3,2,2)`, cell (3,2) rises from 0 to 2, so every prefix covering it grows by 2; the region result becomes **10**.

---

## Key Takeaways

- **2D range sum + point update = 2D Fenwick tree.** Both operations are nested 1D BIT operations → O(log m · log n).
- **Inclusion–exclusion** turns any rectangle sum into four corner-prefix queries — the same identity used for immutable 2D prefix sums (#304), now on an updatable structure.
- **Set vs. add:** cache the grid and push `val − old` so a "set" API works on an additive tree; seeding the tree by Update-ing from zero reuses that same path.
- **Pick by shape:** wide-short grids favor per-row prefixes; square grids with heavy interleaving favor the 2D BIT.

---

## Related Problems

- LeetCode #304 — Range Sum Query 2D - Immutable (2D prefix sums, no updates)
- LeetCode #307 — Range Sum Query - Mutable (1D Fenwick / segment tree)
- LeetCode #303 — Range Sum Query - Immutable (1D prefix sums)
- LeetCode #315 — Count of Smaller Numbers After Self (Fenwick tree on ranks)
