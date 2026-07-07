package main

import "fmt"

// ── Approach 1: Backtracking ──────────────────────────────────────────────────
//
// solveSudoku solves Sudoku Solver using backtracking with constraint tracking.
//
// Intuition: Fill in empty cells one by one. At each empty cell, try digits 1–9.
// If placing a digit is valid (no conflict in row, column, or 3×3 box), recurse.
// If recursion succeeds, done. If not, undo (backtrack) and try the next digit.
// If no digit works, return false to trigger backtracking in the caller.
//
// Algorithm:
//  1. Scan for the first '.' cell.
//  2. For d = '1' to '9':
//     if isValid(board, r, c, d): place d; recurse; if true return true; remove d.
//  3. If no digit fits: return false.
//  4. If no '.' found: board is solved; return true.
//
// Optimisation: Precompute used[row][d], used[col][d], used[box][d] boolean arrays
// instead of scanning the board for validity at each step → O(1) validity check.
//
// Time:  O(9^m) where m = number of empty cells; in practice much faster due to
//        constraint propagation from the precomputed sets.
// Space: O(m) — recursion depth; O(1) for the constraint arrays (fixed 9×9 size)
func solveSudoku(board [][]byte) {
	var rows, cols, boxes [9][10]bool

	// initialise constraint arrays from the given clues
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board[r][c] != '.' {
				d := board[r][c] - '0'
				boxID := (r/3)*3 + (c / 3)
				rows[r][d] = true
				cols[c][d] = true
				boxes[boxID][d] = true
			}
		}
	}

	var backtrack func(pos int) bool
	backtrack = func(pos int) bool {
		// find next empty cell
		for pos < 81 && board[pos/9][pos%9] != '.' {
			pos++
		}
		if pos == 81 { // all cells filled
			return true
		}
		r, c := pos/9, pos%9
		boxID := (r/3)*3 + (c / 3)

		for d := 1; d <= 9; d++ {
			if rows[r][d] || cols[c][d] || boxes[boxID][d] {
				continue // digit already used
			}
			// place digit
			board[r][c] = byte('0' + d)
			rows[r][d] = true
			cols[c][d] = true
			boxes[boxID][d] = true

			if backtrack(pos + 1) {
				return true
			}

			// undo (backtrack)
			board[r][c] = '.'
			rows[r][d] = false
			cols[c][d] = false
			boxes[boxID][d] = false
		}
		return false // no digit worked; trigger backtracking
	}

	backtrack(0)
}

func printBoard(board [][]byte) {
	for r, row := range board {
		if r%3 == 0 && r != 0 {
			fmt.Println("------+-------+------")
		}
		for c, ch := range row {
			if c%3 == 0 && c != 0 {
				fmt.Print("| ")
			}
			fmt.Printf("%c ", ch)
		}
		fmt.Println()
	}
}

func main() {
	board := [][]byte{
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

	fmt.Println("=== Approach 1: Backtracking ===")
	fmt.Println("Input board:")
	printBoard(board)

	solveSudoku(board)

	fmt.Println("\nSolved board:")
	printBoard(board)

	// verify a few known cells from the expected solution
	fmt.Printf("\nboard[0][2]=%c expected '4'\n", board[0][2])
	fmt.Printf("board[1][1]=%c expected '7'\n", board[1][1])
	fmt.Printf("board[8][0]=%c expected '3'\n", board[8][0])
}
