package main

import "fmt"

// ── Approach 1: Group by Diagonal, Reverse Alternate (Brute Force) ────────────
//
// groupByDiagonal solves Diagonal Traverse by first bucketing every cell into
// its anti-diagonal (cells sharing r+c), then emitting the buckets in order,
// reversing every other one so the zig-zag alternates up/down.
//
// Intuition:
//
//	All cells on one anti-diagonal have the same sum r+c. There are m+n-1 such
//	diagonals, indexed d = 0 .. m+n-2. If we collect each diagonal's cells (in
//	natural row order) into a bucket, the required traversal is simply: walk the
//	buckets in increasing d, but read even-numbered diagonals bottom-to-top
//	(reverse) and odd-numbered ones top-to-bottom — that is exactly the
//	up-right / down-left zig-zag LeetCode wants.
//
// Algorithm:
//  1. Bucket cell (r,c) under key d = r+c, appending in increasing r.
//  2. For d = 0 .. m+n-2: if d is even, output bucket[d] reversed; else as-is.
//
// Time:  O(m·n) — every cell is bucketed once and emitted once.
// Space: O(m·n) — the buckets hold all cells before output.
func groupByDiagonal(mat [][]int) []int {
	if len(mat) == 0 || len(mat[0]) == 0 {
		return []int{}
	}
	m, n := len(mat), len(mat[0])
	diagonals := make([][]int, m+n-1) // one bucket per anti-diagonal d = r+c
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			d := r + c // anti-diagonal index; appended in increasing r (top→bottom)
			diagonals[d] = append(diagonals[d], mat[r][c])
		}
	}

	result := make([]int, 0, m*n)
	for d := 0; d < len(diagonals); d++ {
		bucket := diagonals[d]
		if d%2 == 0 {
			// even diagonal → travels UP-right, i.e. bottom row first: reverse it
			for i := len(bucket) - 1; i >= 0; i-- {
				result = append(result, bucket[i])
			}
		} else {
			// odd diagonal → travels DOWN-left, natural top→bottom order
			result = append(result, bucket...)
		}
	}
	return result
}

// ── Approach 2: Simulation with Direction Flips (Optimal, O(1) extra) ─────────
//
// simulateWalk solves Diagonal Traverse by walking the matrix cell-by-cell,
// flipping between the up-right and down-left directions each time it steps off
// an edge, using only O(1) extra space (besides the output).
//
// Intuition:
//
//	Track a current (row, col) and a direction. While going up-right, keep doing
//	r--, c++. When that would leave the grid, "bounce": the next diagonal starts
//	one cell over and the direction flips to down-left (r++, c--). The bounce
//	rules differ by which edge you hit, and the ORDER of the checks matters at
//	the corners (top row vs right column), so handle the right column before the
//	top row when moving up, and the bottom row before the left column when
//	moving down.
//
// Algorithm (moving up-right, direction=+1):
//   - if c == n-1: next start is (r+1, c), flip to down-left  (right edge first)
//   - else if r == 0: next start is (r, c+1), flip to down-left (top edge)
//   - else: r--, c++
//
// Algorithm (moving down-left, direction=-1), mirror:
//   - if r == m-1: next start is (r, c+1), flip to up-right   (bottom edge first)
//   - else if c == 0: next start is (r+1, c), flip to up-right (left edge)
//   - else: r++, c--
//
// Time:  O(m·n) — visits each cell exactly once.
// Space: O(1) — a couple of counters beyond the output slice.
func simulateWalk(mat [][]int) []int {
	if len(mat) == 0 || len(mat[0]) == 0 {
		return []int{}
	}
	m, n := len(mat), len(mat[0])
	result := make([]int, 0, m*n)
	r, c := 0, 0   // current cell, start at top-left
	direction := 1 // +1 = moving up-right, -1 = moving down-left

	for len(result) < m*n { // exactly m*n cells to emit
		result = append(result, mat[r][c]) // record the current cell

		if direction == 1 { // moving up-right (r--, c++)
			switch {
			case c == n-1: // hit right wall → drop down one row, flip to down-left
				r++
				direction = -1
			case r == 0: // hit top wall (and not right wall) → step right, flip
				c++
				direction = -1
			default: // free to keep moving up-right
				r--
				c++
			}
		} else { // moving down-left (r++, c--)
			switch {
			case r == m-1: // hit bottom wall → step right one column, flip to up-right
				c++
				direction = 1
			case c == 0: // hit left wall (and not bottom wall) → drop down, flip
				r++
				direction = 1
			default: // free to keep moving down-left
				r++
				c--
			}
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Group by Diagonal (Brute Force) ===")
	fmt.Println(groupByDiagonal([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}})) // expected [1 2 4 7 5 3 6 8 9]
	fmt.Println(groupByDiagonal([][]int{{1, 2}, {3, 4}}))                  // expected [1 2 3 4]
	fmt.Println(groupByDiagonal([][]int{{1, 2, 3, 4, 5}}))                 // expected [1 2 3 4 5]
	fmt.Println(groupByDiagonal([][]int{{1}, {2}, {3}}))                   // expected [1 2 3]

	fmt.Println("=== Approach 2: Simulation with Direction Flips (Optimal) ===")
	fmt.Println(simulateWalk([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}})) // expected [1 2 4 7 5 3 6 8 9]
	fmt.Println(simulateWalk([][]int{{1, 2}, {3, 4}}))                  // expected [1 2 3 4]
	fmt.Println(simulateWalk([][]int{{1, 2, 3, 4, 5}}))                 // expected [1 2 3 4 5]
	fmt.Println(simulateWalk([][]int{{1}, {2}, {3}}))                   // expected [1 2 3]
}
