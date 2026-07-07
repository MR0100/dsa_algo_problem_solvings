package main

import "fmt"

// ── Approach 1: Extra Matrix ──────────────────────────────────────────────────
//
// extraMatrix solves Rotate Image using a copy of the matrix.
//
// Intuition: When rotating 90° clockwise, element at (r, c) moves to (c, n-1-r).
// Use an extra matrix to hold the result.
//
// Time:  O(n²)
// Space: O(n²) — extra matrix
func extraMatrix(matrix [][]int) {
	n := len(matrix)
	tmp := make([][]int, n)
	for i := range tmp {
		tmp[i] = make([]int, n)
	}
	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			tmp[c][n-1-r] = matrix[r][c] // rotate clockwise
		}
	}
	for r := 0; r < n; r++ {
		copy(matrix[r], tmp[r])
	}
}

// ── Approach 2: Transpose + Reverse Rows (Optimal, In-Place) ─────────────────
//
// rotateInPlace solves Rotate Image in-place using two transformations:
//   1. Transpose: swap matrix[r][c] with matrix[c][r] (reflect across main diagonal).
//   2. Reverse each row (reflect across the vertical midline).
//
// Combined, these two operations produce a 90° clockwise rotation.
//
// Intuition: Rotate clockwise = transpose + reverse rows.
//            Rotate counter-clockwise = transpose + reverse columns.
//
// Time:  O(n²) — n²/2 swaps for transpose + n × n/2 for row reversal
// Space: O(1) — in-place
func rotateInPlace(matrix [][]int) {
	n := len(matrix)

	// step 1: transpose (swap across the main diagonal)
	for r := 0; r < n; r++ {
		for c := r + 1; c < n; c++ { // start from c=r+1 to avoid double-swap
			matrix[r][c], matrix[c][r] = matrix[c][r], matrix[r][c]
		}
	}

	// step 2: reverse each row
	for r := 0; r < n; r++ {
		left, right := 0, n-1
		for left < right {
			matrix[r][left], matrix[r][right] = matrix[r][right], matrix[r][left]
			left++
			right--
		}
	}
}

func printMatrix(m [][]int) {
	for _, row := range m {
		fmt.Println(row)
	}
}

func main() {
	m1 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	m2 := [][]int{{5, 1, 9, 11}, {2, 4, 8, 10}, {13, 3, 6, 7}, {15, 14, 12, 16}}

	fmt.Println("=== Approach 1: Extra Matrix ===")
	a1 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	extraMatrix(a1)
	fmt.Println("Input 3×3 → rotated:")
	printMatrix(a1)
	fmt.Println("Expected: [7 4 1] [8 5 2] [9 6 3]")

	a2 := [][]int{{5, 1, 9, 11}, {2, 4, 8, 10}, {13, 3, 6, 7}, {15, 14, 12, 16}}
	extraMatrix(a2)
	fmt.Println("\nInput 4×4 → rotated:")
	printMatrix(a2)
	fmt.Println("Expected: [15 13 2 5] [14 3 4 1] [12 6 8 9] [16 7 10 11]")

	fmt.Println("\n=== Approach 2: Transpose + Reverse Rows (Optimal) ===")
	b1 := copyMatrix(m1)
	rotateInPlace(b1)
	fmt.Println("Input 3×3 → rotated:")
	printMatrix(b1)
	fmt.Println("Expected: [7 4 1] [8 5 2] [9 6 3]")

	b2 := copyMatrix(m2)
	rotateInPlace(b2)
	fmt.Println("\nInput 4×4 → rotated:")
	printMatrix(b2)
	fmt.Println("Expected: [15 13 2 5] [14 3 4 1] [12 6 8 9] [16 7 10 11]")
}

func copyMatrix(m [][]int) [][]int {
	c := make([][]int, len(m))
	for i := range m {
		c[i] = make([]int, len(m[i]))
		copy(c[i], m[i])
	}
	return c
}
