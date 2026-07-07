package main

import "fmt"

// ── Approach 1: Brute Force (Textbook Triple Loop) ───────────────────────────
//
// bruteForce solves Sparse Matrix Multiplication with the standard definition
// of matrix product: ans[i][j] = sum over p of mat1[i][p] * mat2[p][j].
//
// Intuition:
//
//	The dot product of row i of mat1 with column j of mat2 is exactly the
//	definition of the (i,j) entry of the product. Compute every entry directly.
//	This ignores sparsity entirely and does the full k multiplications per cell.
//
// Algorithm:
//
//  1. Let m = rows of mat1, k = cols of mat1 = rows of mat2, n = cols of mat2.
//  2. For each output cell (i,j), accumulate mat1[i][p]*mat2[p][j] over all p.
//
// Time:  O(m*n*k) — one k-length dot product per output cell.
// Space: O(m*n)   — the result matrix (excluding output).
func bruteForce(mat1 [][]int, mat2 [][]int) [][]int {
	m := len(mat1)    // rows of mat1
	k := len(mat1[0]) // cols of mat1 == rows of mat2 (the shared/contracted dim)
	n := len(mat2[0]) // cols of mat2

	// Allocate the m x n result, zero-initialised.
	ans := make([][]int, m)
	for i := range ans {
		ans[i] = make([]int, n)
	}

	for i := 0; i < m; i++ { // each row of mat1
		for j := 0; j < n; j++ { // each column of mat2
			sum := 0
			for p := 0; p < k; p++ { // dot product over the shared dimension
				sum += mat1[i][p] * mat2[p][j] // multiply even when a factor is 0
			}
			ans[i][j] = sum
		}
	}
	return ans
}

// ── Approach 2: Skip Zeros (Sparsity-Aware Loop Reorder) ──────────────────────
//
// skipZeros solves Sparse Matrix Multiplication by iterating in i-p-j order and
// skipping any mat1[i][p] that is zero, so zero rows/cells cost nothing.
//
// Intuition:
//
//	A zero factor contributes nothing to any output cell. If mat1[i][p] == 0,
//	the entire inner loop over j can be skipped, because every term
//	mat1[i][p]*mat2[p][j] is zero. Reordering loops to i-p-j lets us hoist that
//	check to the p level. For sparse matrices most (i,p) pairs are zero, so we
//	touch only the non-zero products.
//
// Algorithm:
//
//  1. For each row i, for each p: if mat1[i][p] == 0, skip.
//  2. Otherwise, for each j add mat1[i][p]*mat2[p][j] into ans[i][j].
//
// Time:  O(m*k + (#nonzero in mat1)*n) — worst case O(m*n*k), but far less when sparse.
// Space: O(m*n) — the result matrix.
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
				continue // zero factor: whole j-loop would add nothing
			}
			v := mat1[i][p] // hoist the non-zero factor out of the j loop
			for j := 0; j < n; j++ {
				// Only non-zero mat1 entries reach here; still could add 0 if
				// mat2[p][j] is 0, but we avoid the far more common mat1 zeros.
				ans[i][j] += v * mat2[p][j]
			}
		}
	}
	return ans
}

// ── Approach 3: Compressed Sparse Rows (Optimal for very sparse input) ────────
//
// sparseCompressed solves Sparse Matrix Multiplication by first compressing
// each matrix into lists of (column, value) pairs per row, then multiplying
// only non-zero-against-non-zero terms.
//
// Intuition:
//
//	Store, for every row, only its non-zero (col,val) entries. To form the
//	product, for each non-zero mat1[i][p]=v we walk mat2's row p (its non-zero
//	entries only) and scatter v*val into ans[i][col]. Both factors are then
//	guaranteed non-zero, so no work is wasted on zeros at all.
//
// Algorithm:
//
//  1. Build sparse1[i] = list of (p, val) for non-zero mat1[i][p].
//  2. Build sparse2[p] = list of (j, val) for non-zero mat2[p][j].
//  3. For each i, each (p,v1) in sparse1[i], each (j,v2) in sparse2[p]:
//     ans[i][j] += v1*v2.
//
// Time:  O(m*k + k*n + (nonzero pairs actually multiplied)) — proportional to
//
//	real work; near O(nnz1 * avg-nnz-per-row-of-mat2).
//
// Space: O(nnz1 + nnz2 + m*n) — the two sparse representations plus the result.
func sparseCompressed(mat1 [][]int, mat2 [][]int) [][]int {
	m := len(mat1)
	k := len(mat1[0])
	n := len(mat2[0])

	// entry holds a single non-zero cell: its column index and its value.
	type entry struct {
		col int
		val int
	}

	// Compress mat1 by row: sparse1[i] lists non-zero (col, val) of row i.
	sparse1 := make([][]entry, m)
	for i := 0; i < m; i++ {
		for p := 0; p < k; p++ {
			if mat1[i][p] != 0 {
				sparse1[i] = append(sparse1[i], entry{p, mat1[i][p]})
			}
		}
	}

	// Compress mat2 by row: sparse2[p] lists non-zero (col, val) of row p.
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
		for _, e1 := range sparse1[i] { // non-zero factor v1 at column p==e1.col
			p := e1.col
			for _, e2 := range sparse2[p] { // non-zero factor v2 at output col e2.col
				// Both e1.val and e2.val are non-zero: real, useful work only.
				ans[i][e2.col] += e1.val * e2.val
			}
		}
	}
	return ans
}

func main() {
	// Official Example 1
	mat1 := [][]int{{1, 0, 0}, {-1, 0, 3}}
	mat2 := [][]int{{7, 0, 0}, {0, 0, 0}, {0, 0, 1}}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(mat1, mat2)) // expected [[7 0 0] [-7 0 3]]

	fmt.Println("=== Approach 2: Skip Zeros ===")
	fmt.Println(skipZeros(mat1, mat2)) // expected [[7 0 0] [-7 0 3]]

	fmt.Println("=== Approach 3: Compressed Sparse Rows (Optimal) ===")
	fmt.Println(sparseCompressed(mat1, mat2)) // expected [[7 0 0] [-7 0 3]]

	// Official Example 2
	a := [][]int{{0}}
	b := [][]int{{0}}

	fmt.Println("=== Approach 1: Brute Force (Example 2) ===")
	fmt.Println(bruteForce(a, b)) // expected [[0]]

	fmt.Println("=== Approach 2: Skip Zeros (Example 2) ===")
	fmt.Println(skipZeros(a, b)) // expected [[0]]

	fmt.Println("=== Approach 3: Compressed Sparse Rows (Example 2) ===")
	fmt.Println(sparseCompressed(a, b)) // expected [[0]]
}
