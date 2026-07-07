# 0311 — Sparse Matrix Multiplication

> LeetCode #311 · Difficulty: Medium
> **Categories:** Array, Hash Table, Matrix

---

## Problem Statement

Given two sparse matrices `mat1` of size `m x k` and `mat2` of size `k x n`, return the result of `mat1 x mat2`. You may assume that multiplication is always possible.

A **sparse matrix** is a matrix in which most of the elements are zero.

**Example 1:**

```
Input: mat1 = [[1,0,0],[-1,0,3]], mat2 = [[7,0,0],[0,0,0],[0,0,1]]
Output: [[7,0,0],[-7,0,3]]
```

Explanation:
```
mat1 = | 1  0  0 |     mat2 = | 7  0  0 |
       |-1  0  3 |            | 0  0  0 |
                             | 0  0  1 |

result[0][0] = 1*7 + 0*0 + 0*0 = 7
result[1][0] = -1*7 + 0*0 + 3*0 = -7
result[1][2] = -1*0 + 0*0 + 3*1 = 3
```

**Example 2:**

```
Input: mat1 = [[0]], mat2 = [[0]]
Output: [[0]]
```

**Constraints:**

- `m == mat1.length`
- `k == mat1[i].length == mat2.length`
- `n == mat2[i].length`
- `1 <= m, n, k <= 100`
- `-100 <= mat1[i][j], mat2[i][j] <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Matrix traversal** — indexing rows/columns and computing dot products → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **Hash Map / sparse representation** — storing only non-zero cells to skip wasted work → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (triple loop) | O(m·n·k) | O(m·n) | Baseline; small dense matrices |
| 2 | Skip Zeros (loop reorder) | O(m·k + nnz1·n) | O(m·n) | Simple win when input is sparse |
| 3 | Compressed Sparse Rows (Optimal) | O(m·k + k·n + real products) | O(nnz1 + nnz2 + m·n) | Very sparse matrices |

---

## Approach 1 — Brute Force

### Intuition
The `(i,j)` entry of a matrix product is the dot product of row `i` of `mat1` with column `j` of `mat2`. Compute every entry directly by the definition, ignoring sparsity.

### Algorithm
1. Let `m` = rows of `mat1`, `k` = shared dimension, `n` = cols of `mat2`.
2. Allocate an `m x n` zero result.
3. For each cell `(i,j)`, accumulate `mat1[i][p]*mat2[p][j]` over all `p`.

### Complexity
- **Time:** O(m·n·k) — a full `k`-length dot product per output cell.
- **Space:** O(m·n) — the result matrix.

### Code
```go
func bruteForce(mat1 [][]int, mat2 [][]int) [][]int {
	m := len(mat1)
	k := len(mat1[0])
	n := len(mat2[0])

	ans := make([][]int, m)
	for i := range ans {
		ans[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			sum := 0
			for p := 0; p < k; p++ {
				sum += mat1[i][p] * mat2[p][j]
			}
			ans[i][j] = sum
		}
	}
	return ans
}
```

### Dry Run
`mat1 = [[1,0,0],[-1,0,3]]`, `mat2 = [[7,0,0],[0,0,0],[0,0,1]]` (m=2, k=3, n=3).

| i | j | terms `mat1[i][p]*mat2[p][j]` | sum → ans[i][j] |
|---|---|------------------------------|-----------------|
| 0 | 0 | 1·7 + 0·0 + 0·0 | 7 |
| 0 | 1 | 1·0 + 0·0 + 0·0 | 0 |
| 0 | 2 | 1·0 + 0·0 + 0·1 | 0 |
| 1 | 0 | -1·7 + 0·0 + 3·0 | -7 |
| 1 | 1 | -1·0 + 0·0 + 3·0 | 0 |
| 1 | 2 | -1·0 + 0·0 + 3·1 | 3 |

Result: `[[7,0,0],[-7,0,3]]`.

---

## Approach 2 — Skip Zeros

### Intuition
A zero factor contributes nothing to any output cell. If `mat1[i][p] == 0`, the entire inner loop over `j` adds only zeros and can be skipped. Reordering loops to `i-p-j` lets us hoist the zero-check to the `p` level, so common zero entries in `mat1` cost nothing.

### Algorithm
1. For each row `i` and each `p`: if `mat1[i][p] == 0`, `continue`.
2. Otherwise hoist `v = mat1[i][p]` and for each `j` add `v*mat2[p][j]` into `ans[i][j]`.

### Complexity
- **Time:** O(m·k + nnz1·n) — where `nnz1` is the number of non-zeros in `mat1`; worst case O(m·n·k) when dense.
- **Space:** O(m·n) — the result matrix.

### Code
```go
func skipZeros(mat1 [][]int, mat2 [][]int) [][]int {
	m := len(mat1)
	k := len(mat1[0])
	n := len(mat2[0])

	ans := make([][]int, m)
	for i := range ans {
		ans[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		for p := 0; p < k; p++ {
			if mat1[i][p] == 0 {
				continue
			}
			v := mat1[i][p]
			for j := 0; j < n; j++ {
				ans[i][j] += v * mat2[p][j]
			}
		}
	}
	return ans
}
```

### Dry Run
`mat1 = [[1,0,0],[-1,0,3]]`, `mat2 = [[7,0,0],[0,0,0],[0,0,1]]`.

| i | p | mat1[i][p] | action | ans row after |
|---|---|-----------|--------|---------------|
| 0 | 0 | 1 | v=1, add 1·[7,0,0] | ans[0]=[7,0,0] |
| 0 | 1 | 0 | skip | ans[0]=[7,0,0] |
| 0 | 2 | 0 | skip | ans[0]=[7,0,0] |
| 1 | 0 | -1 | v=-1, add -1·[7,0,0] | ans[1]=[-7,0,0] |
| 1 | 1 | 0 | skip | ans[1]=[-7,0,0] |
| 1 | 2 | 3 | v=3, add 3·[0,0,1] | ans[1]=[-7,0,3] |

Result: `[[7,0,0],[-7,0,3]]`.

---

## Approach 3 — Compressed Sparse Rows (Optimal)

### Intuition
Store, per row, only the non-zero `(col, val)` entries of each matrix. To multiply, for each non-zero `mat1[i][p]=v1` walk only `mat2`'s row `p` non-zero entries `(j, v2)` and scatter `v1*v2` into `ans[i][j]`. Both factors are guaranteed non-zero, so zero work is wasted.

### Algorithm
1. Build `sparse1[i]` = list of `(p, val)` for non-zero `mat1[i][p]`.
2. Build `sparse2[p]` = list of `(j, val)` for non-zero `mat2[p][j]`.
3. For each `i`, each `(p, v1)` in `sparse1[i]`, each `(j, v2)` in `sparse2[p]`: `ans[i][j] += v1*v2`.

### Complexity
- **Time:** O(m·k + k·n + real products) — proportional to actual non-zero multiplications.
- **Space:** O(nnz1 + nnz2 + m·n) — the two sparse representations plus the result.

### Code
```go
func sparseCompressed(mat1 [][]int, mat2 [][]int) [][]int {
	m := len(mat1)
	k := len(mat1[0])
	n := len(mat2[0])

	type entry struct {
		col int
		val int
	}

	sparse1 := make([][]entry, m)
	for i := 0; i < m; i++ {
		for p := 0; p < k; p++ {
			if mat1[i][p] != 0 {
				sparse1[i] = append(sparse1[i], entry{p, mat1[i][p]})
			}
		}
	}

	sparse2 := make([][]entry, k)
	for p := 0; p < k; p++ {
		for j := 0; j < n; j++ {
			if mat2[p][j] != 0 {
				sparse2[p] = append(sparse2[p], entry{j, mat2[p][j]})
			}
		}
	}

	ans := make([][]int, m)
	for i := range ans {
		ans[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		for _, e1 := range sparse1[i] {
			p := e1.col
			for _, e2 := range sparse2[p] {
				ans[i][e2.col] += e1.val * e2.val
			}
		}
	}
	return ans
}
```

### Dry Run
`mat1 = [[1,0,0],[-1,0,3]]`, `mat2 = [[7,0,0],[0,0,0],[0,0,1]]`.

Compress: `sparse1 = [[(0,1)], [(0,-1),(2,3)]]`, `sparse2 = [[(0,7)], [], [(2,1)]]`.

| i | e1 (p,v1) | e2 (j,v2) in sparse2[p] | update |
|---|-----------|--------------------------|--------|
| 0 | (0,1) | (0,7) | ans[0][0] += 1·7 = 7 |
| 1 | (0,-1) | (0,7) | ans[1][0] += -1·7 = -7 |
| 1 | (2,3) | (2,1) | ans[1][2] += 3·1 = 3 |

Result: `[[7,0,0],[-7,0,3]]`. Note only 3 multiplications happened vs 18 in brute force.

---

## Key Takeaways
- **Loop reordering** (`i-p-j` instead of `i-j-p`) exposes a natural spot to skip zeros, turning dense matrix multiply into sparse-aware multiply with a one-line `continue`.
- **Compressed sparse representation** (list of non-zero `(col,val)` per row) is the standard trick for sparse linear algebra; it makes cost proportional to the number of real products.
- The `k` shared dimension is both the number of columns of `mat1` and rows of `mat2`; always index the contracted dimension consistently.

---

## Related Problems
- LeetCode #48 — Rotate Image (matrix traversal)
- LeetCode #54 — Spiral Matrix (matrix traversal)
- LeetCode #73 — Set Matrix Zeroes (sparse-cell bookkeeping)
- LeetCode #1 — Two Sum (hash map of seen values)
