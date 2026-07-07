package main

import "fmt"

// ── Approach 1: Backtracking (Count Only) ────────────────────────────────────
//
// backtracking solves N-Queens II: count distinct solutions without storing boards.
//
// Intuition:
//   Same as #51 but instead of recording the board, just increment a counter
//   when all n rows are placed. Skip storing or reconstructing the board.
//
// Algorithm:
//   Place one queen per row; use cols/diag/anti sets to check conflicts.
//   Increment count when row == n.
//
// Time:  O(n!)
// Space: O(n)
func backtracking(n int) int {
	count := 0
	cols := make([]bool, n)
	diag := make([]bool, 2*n)
	anti := make([]bool, 2*n)

	var bt func(row int)
	bt = func(row int) {
		if row == n {
			count++
			return
		}
		for c := 0; c < n; c++ {
			dIdx := (row-c) + (n - 1)
			aIdx := row + c
			if cols[c] || diag[dIdx] || anti[aIdx] {
				continue
			}
			cols[c] = true
			diag[dIdx] = true
			anti[aIdx] = true

			bt(row + 1)

			cols[c] = false
			diag[dIdx] = false
			anti[aIdx] = false
		}
	}

	bt(0)
	return count
}

// ── Approach 2: Bitmask Backtracking ─────────────────────────────────────────
//
// bitmask solves N-Queens II using integer bitmasks for O(1) bit operations.
//
// Intuition:
//   Encode which columns are attacked as a bitmask `cols`. The occupied
//   diagonals propagate: when we move to row+1, left-diagonals shift right
//   by 1 (>>1) and right-diagonals shift left (<<1). The available columns
//   for the next row are those not set in cols | leftDiag | rightDiag.
//
//   Iterate set bits of `avail` using `avail & (-avail)` (lowest set bit),
//   place a queen there, then recurse with updated masks.
//
// Algorithm:
//   bt(row, cols, leftDiag, rightDiag):
//     avail = ALL_MASK & ^(cols | leftDiag | rightDiag)
//     while avail != 0:
//       bit = avail & (-avail)  // pick lowest available column
//       avail &= avail-1         // clear that bit
//       bt(row+1, cols|bit, (leftDiag|bit)>>1, (rightDiag|bit)<<1)
//
// Time:  O(n!)  — same recursion tree, but constant-time bit ops per level.
// Space: O(n)   — recursion stack only; no boolean arrays.
func bitmask(n int) int {
	count := 0
	allMask := (1 << n) - 1 // n bits all set: valid column range

	var bt func(row, cols, leftDiag, rightDiag int)
	bt = func(row, cols, leftDiag, rightDiag int) {
		if row == n {
			count++
			return
		}
		// bits that are NOT attacked in this row
		avail := allMask & ^(cols | leftDiag | rightDiag)
		for avail != 0 {
			bit := avail & (-avail) // isolate lowest set bit
			avail &= avail - 1      // remove that bit from candidates
			// diagonals shift one step each row
			bt(row+1, cols|bit, (leftDiag|bit)>>1, (rightDiag|bit)<<1)
		}
	}

	bt(0, 0, 0, 0)
	return count
}

func main() {
	fmt.Println("=== Approach 1: Backtracking ===")
	fmt.Printf("n=4  count=%d  expected 2\n", backtracking(4))
	fmt.Printf("n=1  count=%d  expected 1\n", backtracking(1))
	fmt.Printf("n=8  count=%d  expected 92\n", backtracking(8))

	fmt.Println("=== Approach 2: Bitmask Backtracking ===")
	fmt.Printf("n=4  count=%d  expected 2\n", bitmask(4))
	fmt.Printf("n=1  count=%d  expected 1\n", bitmask(1))
	fmt.Printf("n=8  count=%d  expected 92\n", bitmask(8))
}
