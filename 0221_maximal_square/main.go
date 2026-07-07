package main

import "fmt"

// ── Approach 1: Brute Force (Expand Every Square) ────────────────────────────
//
// bruteForce solves Maximal Square by, from every cell that holds a '1',
// trying to grow the largest all-ones square whose top-left corner sits there.
//
// Intuition:
//
//	A square of side k with top-left (r,c) is all ones iff every one of its
//	k² cells is '1'. So for each '1' cell we try side = 1, 2, 3, … and stop
//	the moment we hit a '0' or run off the grid; the largest side that still
//	worked, over all cells, gives the answer (squared for the area).
//
// Algorithm:
//
//  1. For each cell (r,c) that is '1':
//     - Try to extend the current best+1 sized square (only bigger helps).
//     - For a candidate side, scan the new right column and bottom row of
//     cells the square would add; if all are '1', accept the larger side.
//  2. Track the maximum side seen; answer is side².
//
// Time:  O(m·n·min(m,n)²) worst case — each cell may re-scan an O(side²) block.
// Space: O(1) — only scalar bookkeeping.
func bruteForce(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0 // empty grid → no square
	}
	m, n := len(matrix), len(matrix[0])
	best := 0 // largest side length found so far
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if matrix[r][c] != '1' {
				continue // a square's top-left must itself be '1'
			}
			side := 1 // a lone '1' is already a 1×1 square
			// Try to grow the square one ring at a time.
			for r+side < m && c+side < n && valid(matrix, r, c, side) {
				side++ // the (side+1)×(side+1) square is all ones too
			}
			if side > best {
				best = side // remember the biggest square so far
			}
		}
	}
	return best * best // area = side²
}

// valid reports whether extending the square at (r,c) from size `side` to
// `side+1` keeps it all ones — i.e. the freshly-added L-shaped border
// (new bottom row + new right column) is entirely '1'.
func valid(matrix [][]byte, r, c, side int) bool {
	for i := 0; i <= side; i++ {
		if matrix[r+side][c+i] != '1' { // new bottom row
			return false
		}
		if matrix[r+i][c+side] != '1' { // new right column
			return false
		}
	}
	return true
}

// ── Approach 2: Dynamic Programming (2D Table) ───────────────────────────────
//
// dp2D solves Maximal Square with the classic DP where dp[r][c] is the side of
// the largest all-ones square whose *bottom-right* corner is (r,c).
//
// Intuition:
//
//	A square of side k ending at (r,c) requires three overlapping squares of
//	side k-1 ending at its left, top, and top-left neighbours to all exist —
//	their minimum limits how far this corner can grow. Hence
//	dp[r][c] = 1 + min(dp[r-1][c], dp[r][c-1], dp[r-1][c-1]) when the cell is
//	'1', else 0. The bottleneck (the smallest neighbour) is what caps growth.
//
// Algorithm:
//
//  1. dp[r][c] = 0 for every '0' cell.
//  2. For a '1' cell on the top row or left column, dp = 1 (no room to grow).
//  3. Otherwise dp[r][c] = 1 + min of the three neighbours above/left/diagonal.
//  4. Track the max dp value; answer is that side².
//
// Time:  O(m·n) — one O(1) recurrence per cell.
// Space: O(m·n) — the full dp table.
func dp2D(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])
	dp := make([][]int, m) // dp[r][c] = largest square side ending at (r,c)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	best := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if matrix[r][c] != '1' {
				continue // '0' cell ends no square → dp stays 0
			}
			if r == 0 || c == 0 {
				dp[r][c] = 1 // border cell: at most a 1×1 square
			} else {
				// limited by the smallest of the three overlapping squares
				dp[r][c] = 1 + min3(dp[r-1][c], dp[r][c-1], dp[r-1][c-1])
			}
			if dp[r][c] > best {
				best = dp[r][c] // track the biggest square side
			}
		}
	}
	return best * best
}

// ── Approach 3: Dynamic Programming (1D Rolling Row) — Optimal ───────────────
//
// dp1D solves Maximal Square using the same recurrence but keeping only one
// rolling row of dp values plus a single saved "top-left" scalar.
//
// Intuition:
//
//	dp[r][c] only ever reads row r and row r-1, so a full 2D table is wasted
//	memory. Overwrite a single array left-to-right; before overwriting dp[c]
//	stash its old value (that is dp[r-1][c-1] for the next column) in `prev`.
//	This drops space from O(m·n) to O(n) with identical results.
//
// Algorithm:
//
//  1. Keep a 1D slice `dp` of length n+1 (index shifted by 1 to avoid a
//     special case for the first column) initialised to 0.
//  2. Walk rows top→bottom, columns left→right. Hold `prev` = dp[c] before it
//     is overwritten (the diagonal top-left value).
//  3. For a '1' cell: dp[c] = 1 + min(dp[c] (top), dp[c-1] (left), prev
//     (top-left)); for a '0' cell dp[c] = 0.
//  4. Track the max; answer is side².
//
// Time:  O(m·n) — one pass, O(1) per cell.
// Space: O(n) — a single rolling row.
func dp1D(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])
	dp := make([]int, n+1) // dp[c+1] = square side ending at current row, column c
	best := 0
	for r := 0; r < m; r++ {
		prev := 0 // dp[r-1][c-1]: top-left diagonal, starts 0 each row
		for c := 0; c < n; c++ {
			temp := dp[c+1] // save dp[r-1][c] before overwriting (becomes next prev)
			if matrix[r][c] == '1' {
				// dp[c+1]=top, dp[c]=left, prev=top-left
				dp[c+1] = 1 + min3(dp[c+1], dp[c], prev)
				if dp[c+1] > best {
					best = dp[c+1]
				}
			} else {
				dp[c+1] = 0 // '0' cell resets the running square
			}
			prev = temp // this row's dp[c] is next column's top-left
		}
	}
	return best * best
}

// min3 returns the smallest of three ints.
func min3(a, b, c int) int {
	if b < a {
		a = b
	}
	if c < a {
		a = c
	}
	return a
}

func main() {
	// Example 1
	m1 := [][]byte{
		{'1', '0', '1', '0', '0'},
		{'1', '0', '1', '1', '1'},
		{'1', '1', '1', '1', '1'},
		{'1', '0', '0', '1', '0'},
	}
	// Example 2
	m2 := [][]byte{
		{'0', '1'},
		{'1', '0'},
	}
	// Example 3
	m3 := [][]byte{
		{'0'},
	}

	fmt.Println("=== Approach 1: Brute Force (Expand Every Square) ===")
	fmt.Println(bruteForce(cloneGrid(m1))) // expected 4
	fmt.Println(bruteForce(cloneGrid(m2))) // expected 1
	fmt.Println(bruteForce(cloneGrid(m3))) // expected 0

	fmt.Println("=== Approach 2: Dynamic Programming (2D Table) ===")
	fmt.Println(dp2D(cloneGrid(m1))) // expected 4
	fmt.Println(dp2D(cloneGrid(m2))) // expected 1
	fmt.Println(dp2D(cloneGrid(m3))) // expected 0

	fmt.Println("=== Approach 3: Dynamic Programming (1D Rolling Row) (Optimal) ===")
	fmt.Println(dp1D(cloneGrid(m1))) // expected 4
	fmt.Println(dp1D(cloneGrid(m2))) // expected 1
	fmt.Println(dp1D(cloneGrid(m3))) // expected 0
}

// cloneGrid deep-copies a byte grid so each approach gets a fresh, unmodified
// input (none of these mutate, but this keeps the driver defensive and clear).
func cloneGrid(g [][]byte) [][]byte {
	out := make([][]byte, len(g))
	for i := range g {
		out[i] = append([]byte(nil), g[i]...)
	}
	return out
}
