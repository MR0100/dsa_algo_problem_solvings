package main

import "fmt"

// ── Approach 1: Backtracking with Visited Sets ───────────────────────────────
//
// backtracking solves N-Queens using row-by-row backtracking with three
// boolean sets tracking which columns, diagonals, and anti-diagonals are occupied.
//
// Intuition:
//   Place one queen per row. At each row, try every column that isn't attacked
//   by an existing queen. A queen at (r,c) attacks column c, diagonal r-c,
//   and anti-diagonal r+c.
//
// Algorithm:
//   1. For each row r (0..n-1), iterate columns c (0..n-1).
//   2. Skip if cols[c], diag[r-c], or anti[r+c] is true.
//   3. Mark all three, recurse to row r+1.
//   4. If r==n, convert the board to string representation and record.
//   5. Unmark on backtrack.
//
// Time:  O(n!) — at most n! placements in the worst case.
// Space: O(n)  — recursion depth n; sets are O(n).
func backtracking(n int) [][]string {
	var result [][]string
	cols := make(map[int]bool)
	diag := make(map[int]bool)  // r - c is constant on each diagonal
	anti := make(map[int]bool)  // r + c is constant on each anti-diagonal
	board := make([][]byte, n)
	for i := range board {
		board[i] = make([]byte, n)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}

	var bt func(row int)
	bt = func(row int) {
		if row == n {
			// convert board snapshot to []string
			snap := make([]string, n)
			for i, r := range board {
				snap[i] = string(r)
			}
			result = append(result, snap)
			return
		}
		for c := 0; c < n; c++ {
			if cols[c] || diag[row-c] || anti[row+c] {
				continue // column or diagonal already attacked
			}
			// place queen
			board[row][c] = 'Q'
			cols[c] = true
			diag[row-c] = true
			anti[row+c] = true

			bt(row + 1)

			// remove queen (backtrack)
			board[row][c] = '.'
			cols[c] = false
			diag[row-c] = false
			anti[row+c] = false
		}
	}

	bt(0)
	return result
}

// ── Approach 2: Backtracking with Boolean Arrays ──────────────────────────────
//
// backtrackingArrays solves N-Queens using arrays instead of maps for the
// attacked-set lookups, which is faster in practice due to array indexing.
//
// Intuition:
//   Same algorithm as Approach 1. Diagonals r-c range from -(n-1) to n-1,
//   so offset by n-1 to index into a boolean array. Anti-diagonals r+c range
//   from 0 to 2*(n-1), so no offset needed.
//
// Algorithm:
//   Identical to Approach 1 but cols/diag/anti are []bool slices instead of maps.
//
// Time:  O(n!)
// Space: O(n)
func backtrackingArrays(n int) [][]string {
	var result [][]string
	cols := make([]bool, n)
	diag := make([]bool, 2*n)  // index: (r-c) + (n-1) to keep non-negative
	anti := make([]bool, 2*n)  // index: r+c
	board := make([][]byte, n)
	for i := range board {
		board[i] = make([]byte, n)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}

	var bt func(row int)
	bt = func(row int) {
		if row == n {
			snap := make([]string, n)
			for i, r := range board {
				snap[i] = string(r)
			}
			result = append(result, snap)
			return
		}
		for c := 0; c < n; c++ {
			dIdx := (row - c) + (n - 1) // shift to keep non-negative
			aIdx := row + c
			if cols[c] || diag[dIdx] || anti[aIdx] {
				continue
			}
			board[row][c] = 'Q'
			cols[c] = true
			diag[dIdx] = true
			anti[aIdx] = true

			bt(row + 1)

			board[row][c] = '.'
			cols[c] = false
			diag[dIdx] = false
			anti[aIdx] = false
		}
	}

	bt(0)
	return result
}

func main() {
	fmt.Println("=== Approach 1: Backtracking with Maps ===")
	r1 := backtracking(4)
	fmt.Printf("n=4  solutions=%d  expected 2\n", len(r1))
	for _, sol := range r1 {
		for _, row := range sol {
			fmt.Println(" ", row)
		}
		fmt.Println()
	}

	r2 := backtracking(1)
	fmt.Printf("n=1  solutions=%d  expected 1\n", len(r2))

	fmt.Println("=== Approach 2: Backtracking with Arrays ===")
	r3 := backtrackingArrays(4)
	fmt.Printf("n=4  solutions=%d  expected 2\n", len(r3))
	for _, sol := range r3 {
		for _, row := range sol {
			fmt.Println(" ", row)
		}
		fmt.Println()
	}

	r4 := backtrackingArrays(1)
	fmt.Printf("n=1  solutions=%d  expected 1\n", len(r4))
}
