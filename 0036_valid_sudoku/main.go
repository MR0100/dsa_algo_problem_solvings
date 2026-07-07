package main

import "fmt"

// ── Approach 1: Three Separate Passes ────────────────────────────────────────
//
// threePasses solves Valid Sudoku by checking rows, columns, and 3×3 boxes in
// three separate passes over the board.
//
// Intuition: Validity requires no digit repeats within any row, column, or box.
// Separate passes make the logic easiest to read.
//
// Time:  O(1) — board is always 9×9 = 81 cells; bounded constant
// Space: O(1) — at most 9 booleans per pass; bounded constant
func threePasses(board [][]byte) bool {
	// check all rows
	for r := 0; r < 9; r++ {
		seen := [10]bool{}
		for c := 0; c < 9; c++ {
			if board[r][c] == '.' {
				continue
			}
			d := board[r][c] - '0'
			if seen[d] {
				return false
			}
			seen[d] = true
		}
	}
	// check all columns
	for c := 0; c < 9; c++ {
		seen := [10]bool{}
		for r := 0; r < 9; r++ {
			if board[r][c] == '.' {
				continue
			}
			d := board[r][c] - '0'
			if seen[d] {
				return false
			}
			seen[d] = true
		}
	}
	// check all 3×3 boxes
	for boxRow := 0; boxRow < 3; boxRow++ {
		for boxCol := 0; boxCol < 3; boxCol++ {
			seen := [10]bool{}
			for r := boxRow * 3; r < boxRow*3+3; r++ {
				for c := boxCol * 3; c < boxCol*3+3; c++ {
					if board[r][c] == '.' {
						continue
					}
					d := board[r][c] - '0'
					if seen[d] {
						return false
					}
					seen[d] = true
				}
			}
		}
	}
	return true
}

// ── Approach 2: Single Pass with Three Seen Arrays (Optimal) ─────────────────
//
// singlePass solves Valid Sudoku in a single pass over all 81 cells, checking
// the row, column, and box constraints simultaneously.
//
// Intuition: For each cell (r, c) with digit d:
//   - rows[r][d]: has digit d appeared in row r?
//   - cols[c][d]: has digit d appeared in column c?
//   - boxes[boxId][d]: has digit d appeared in box boxId?
//   boxId = (r/3)*3 + (c/3) maps each cell to its 3×3 box index (0–8).
//
// Time:  O(1) — 81 cells × constant work per cell
// Space: O(1) — fixed-size arrays (9×9 each for rows/cols/boxes)
func singlePass(board [][]byte) bool {
	var rows, cols, boxes [9][10]bool

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board[r][c] == '.' {
				continue
			}
			d := board[r][c] - '0'           // digit 1–9
			boxID := (r/3)*3 + (c / 3)       // box index 0–8

			if rows[r][d] || cols[c][d] || boxes[boxID][d] {
				return false // duplicate found
			}
			rows[r][d] = true
			cols[c][d] = true
			boxes[boxID][d] = true
		}
	}
	return true
}

func main() {
	valid := [][]byte{
		{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
		{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
		{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
		{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
		{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
		{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
		{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
		{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
		{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
	}

	invalid := [][]byte{
		{'8', '3', '.', '.', '7', '.', '.', '.', '.'},
		{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
		{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
		{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
		{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
		{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
		{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
		{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
		{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
	}

	fmt.Println("=== Approach 1: Three Passes ===")
	fmt.Printf("valid board:   %v  expected true\n", threePasses(valid))
	fmt.Printf("invalid board: %v  expected false\n", threePasses(invalid))

	fmt.Println("\n=== Approach 2: Single Pass (Optimal) ===")
	fmt.Printf("valid board:   %v  expected true\n", singlePass(valid))
	fmt.Printf("invalid board: %v  expected false\n", singlePass(invalid))
}
