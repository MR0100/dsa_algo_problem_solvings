package main

import "fmt"

// ── Approach 1: Extra Space ───────────────────────────────────────────────────
//
// extraSpace solves Set Matrix Zeroes using separate row and column flag arrays.
//
// Intuition:
//   First pass: record which rows and columns contain zeros.
//   Second pass: for each cell, if its row or column is flagged, set to 0.
//
// Time:  O(m × n)
// Space: O(m + n)
func extraSpace(matrix [][]int) {
	m, n := len(matrix), len(matrix[0])
	rows := make([]bool, m)
	cols := make([]bool, n)

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if matrix[r][c] == 0 {
				rows[r] = true
				cols[c] = true
			}
		}
	}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if rows[r] || cols[c] {
				matrix[r][c] = 0
			}
		}
	}
}

// ── Approach 2: In-Place with First Row/Col as Flags ─────────────────────────
//
// inPlace solves Set Matrix Zeroes in O(1) extra space by using the first row
// and first column of the matrix itself as flag arrays.
//
// Intuition:
//   Use matrix[0][c] to flag column c, and matrix[r][0] to flag row r.
//   The cell matrix[0][0] is shared by both row 0 and col 0, so we need a
//   separate boolean (firstColZero) to track whether column 0 was originally zero.
//
// Algorithm:
//   1. Check if first row/column has any zero (save as flags).
//   2. Mark flags in first row/col based on interior cells.
//   3. Zero out interior cells based on first row/col flags.
//   4. Zero out first row/col if their respective flags are set.
//
// Time:  O(m × n)
// Space: O(1)
func inPlace(matrix [][]int) {
	m, n := len(matrix), len(matrix[0])

	// does column 0 need to be zeroed?
	firstColZero := false
	for r := 0; r < m; r++ {
		if matrix[r][0] == 0 {
			firstColZero = true
			break
		}
	}

	// does row 0 need to be zeroed?
	firstRowZero := false
	for c := 0; c < n; c++ {
		if matrix[0][c] == 0 {
			firstRowZero = true
			break
		}
	}

	// use first row and first column as flags for the rest of the matrix
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			if matrix[r][c] == 0 {
				matrix[r][0] = 0 // flag row r
				matrix[0][c] = 0 // flag col c
			}
		}
	}

	// zero interior cells based on flags
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			if matrix[r][0] == 0 || matrix[0][c] == 0 {
				matrix[r][c] = 0
			}
		}
	}

	// zero first row if needed
	if firstRowZero {
		for c := 0; c < n; c++ {
			matrix[0][c] = 0
		}
	}

	// zero first column if needed
	if firstColZero {
		for r := 0; r < m; r++ {
			matrix[r][0] = 0
		}
	}
}

func printMatrix(m [][]int) {
	for _, row := range m {
		fmt.Println(" ", row)
	}
}

func main() {
	fmt.Println("=== Approach 1: Extra Space ===")
	m1 := [][]int{{1, 1, 1}, {1, 0, 1}, {1, 1, 1}}
	extraSpace(m1)
	fmt.Println("Input [[1,1,1],[1,0,1],[1,1,1]] got:")
	printMatrix(m1)
	fmt.Println("Expected: [[1,0,1],[0,0,0],[1,0,1]]")

	m2 := [][]int{{0, 1, 2, 0}, {3, 4, 5, 2}, {1, 3, 1, 5}}
	extraSpace(m2)
	fmt.Println("Input [[0,1,2,0],[3,4,5,2],[1,3,1,5]] got:")
	printMatrix(m2)
	fmt.Println("Expected: [[0,0,0,0],[0,4,5,0],[0,3,1,0]]")

	fmt.Println("=== Approach 2: In-Place ===")
	m3 := [][]int{{1, 1, 1}, {1, 0, 1}, {1, 1, 1}}
	inPlace(m3)
	fmt.Println("Input [[1,1,1],[1,0,1],[1,1,1]] got:")
	printMatrix(m3)

	m4 := [][]int{{0, 1, 2, 0}, {3, 4, 5, 2}, {1, 3, 1, 5}}
	inPlace(m4)
	fmt.Println("Input [[0,1,2,0],[3,4,5,2],[1,3,1,5]] got:")
	printMatrix(m4)
}
