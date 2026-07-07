package main

import "fmt"

// LeetCode 304 — Range Sum Query 2D - Immutable
//
// Given a 2D matrix, answer many queries for the sum of the sub-rectangle
// bounded by (row1, col1) top-left and (row2, col2) bottom-right, inclusive.
// The matrix is immutable, so preprocess once and answer each query in O(1).

// ── Approach 1: Brute Force (Re-sum Each Query) ──────────────────────────────
//
// NumMatrixBrute keeps the raw matrix and re-adds the sub-rectangle each query.
//
// Intuition:
//
//	Simplest possible: store the matrix and, for each query, double-loop over
//	the requested rows and columns adding every cell. Correct, but O(rows·cols)
//	per query — far too slow when there are many queries.
//
// Time:  constructor O(1); SumRegion O(rows·cols) per query.
// Space: O(m·n) — stores the matrix.
type NumMatrixBrute struct {
	matrix [][]int // the original, immutable grid
}

// NewNumMatrixBrute stores the grid by reference (never mutated).
func NewNumMatrixBrute(matrix [][]int) NumMatrixBrute {
	return NumMatrixBrute{matrix: matrix}
}

// SumRegion adds every cell in the requested rectangle on demand.
//
// Time:  O((row2−row1+1)·(col2−col1+1)).
// Space: O(1).
func (m NumMatrixBrute) SumRegion(row1, col1, row2, col2 int) int {
	sum := 0
	for r := row1; r <= row2; r++ { // each row of the rectangle
		for c := col1; c <= col2; c++ { // each column of the rectangle
			sum += m.matrix[r][c]
		}
	}
	return sum
}

// ── Approach 2: 2D Prefix Sum (Optimal) ──────────────────────────────────────
//
// NumMatrix precomputes cumulative area sums so any sub-rectangle answers O(1).
//
// Intuition:
//
//	Define prefix[r][c] = sum of the whole sub-rectangle from (0,0) to
//	(r−1, c−1). Any rectangle sum then comes from inclusion–exclusion:
//
//	    sum = P[r2+1][c2+1]   (whole area up to bottom-right)
//	        − P[r1][c2+1]     (strip above the rectangle)
//	        − P[r2+1][c1]     (strip left of the rectangle)
//	        + P[r1][c1]       (top-left corner added back — subtracted twice)
//
//	The +1 padding (an all-zero first row and column) removes boundary special
//	cases. One O(m·n) build gives O(1) queries.
//
// Algorithm:
//  1. Build prefix of size (m+1)×(n+1), all-zero border, where
//     prefix[r+1][c+1] = matrix[r][c] + prefix[r][c+1] + prefix[r+1][c] − prefix[r][c].
//  2. SumRegion uses the four-corner inclusion–exclusion formula above.
//
// Time:  constructor O(m·n); SumRegion O(1) per query.
// Space: O(m·n) — the prefix grid.
type NumMatrix struct {
	prefix [][]int // (m+1)×(n+1) padded cumulative-area sums
}

// NewNumMatrix builds the 2D cumulative-sum table once.
//
// Time:  O(m·n).
// Space: O(m·n).
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
			// Cumulative area up to (r,c): current cell + area above + area
			// to the left − the doubly-counted top-left overlap.
			prefix[r+1][c+1] = matrix[r][c] +
				prefix[r][c+1] +
				prefix[r+1][c] -
				prefix[r][c]
		}
	}
	return NumMatrix{prefix: prefix}
}

// SumRegion returns the rectangle sum via four-corner inclusion–exclusion.
//
// Time:  O(1).
// Space: O(1).
func (nm NumMatrix) SumRegion(row1, col1, row2, col2 int) int {
	p := nm.prefix
	return p[row2+1][col2+1] - // full area to bottom-right corner
		p[row1][col2+1] - // remove the strip above the rectangle
		p[row2+1][col1] + // remove the strip to the left
		p[row1][col1] //     add back the top-left, subtracted twice
}

func main() {
	matrix := [][]int{
		{3, 0, 1, 4, 2},
		{5, 6, 3, 2, 1},
		{1, 2, 0, 1, 5},
		{4, 1, 0, 1, 7},
		{1, 0, 3, 0, 5},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	mb := NewNumMatrixBrute(matrix)
	fmt.Println(mb.SumRegion(2, 1, 4, 3)) // expected 8
	fmt.Println(mb.SumRegion(1, 1, 2, 2)) // expected 11
	fmt.Println(mb.SumRegion(1, 2, 2, 4)) // expected 12

	fmt.Println("=== Approach 2: 2D Prefix Sum (Optimal) ===")
	nm := NewNumMatrix(matrix)
	fmt.Println(nm.SumRegion(2, 1, 4, 3)) // expected 8
	fmt.Println(nm.SumRegion(1, 1, 2, 2)) // expected 11
	fmt.Println(nm.SumRegion(1, 2, 2, 4)) // expected 12
}
