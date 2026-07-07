package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Search a 2D Matrix II by scanning every cell.
//
// Intuition:
//
//	Ignore the sorted structure entirely and check each element. Obviously
//	correct, the baseline against which the smarter methods are measured.
//
// Algorithm:
//  1. Loop over every row and column.
//  2. Return true if any cell equals target.
//  3. Return false if the scan completes with no match.
//
// Time:  O(m·n) — visits every cell.
// Space: O(1).
func bruteForce(matrix [][]int, target int) bool {
	for _, row := range matrix {
		for _, v := range row {
			if v == target {
				return true // found it
			}
		}
	}
	return false // no cell matched
}

// ── Approach 2: Binary Search Each Row ───────────────────────────────────────
//
// binarySearchRows solves Search a 2D Matrix II by binary-searching within each
// (individually sorted) row.
//
// Intuition:
//
//	Every row is sorted left→right, so a per-row binary search finds the target
//	in O(log n). Doing that for all m rows is O(m log n) — better than brute
//	force, though it does not exploit the column ordering.
//
// Algorithm:
//  1. For each row, binary-search for target.
//  2. Return true on the first hit; false if none.
//
// Time:  O(m log n) — m binary searches over n columns.
// Space: O(1).
func binarySearchRows(matrix [][]int, target int) bool {
	for _, row := range matrix {
		lo, hi := 0, len(row)-1
		for lo <= hi {
			mid := lo + (hi-lo)/2 // avoid overflow
			switch {
			case row[mid] == target:
				return true
			case row[mid] < target:
				lo = mid + 1 // target is to the right
			default:
				hi = mid - 1 // target is to the left
			}
		}
	}
	return false
}

// ── Approach 3: Staircase Search from Top-Right (Optimal) ─────────────────────
//
// staircaseSearch solves Search a 2D Matrix II by starting at the top-right
// corner and eliminating one full row or column per step.
//
// Intuition:
//
//	Stand at the top-right cell. It is the largest in its row and the smallest
//	in its column. So if it's bigger than target, nothing below it in this
//	column can help → move LEFT (drop the column). If it's smaller than target,
//	nothing to its left in this row can help → move DOWN (drop the row). Each
//	step removes an entire row or column, giving O(m+n).
//
// Algorithm:
//  1. Start at row = 0, col = n-1.
//  2. While in bounds: compare matrix[row][col] to target.
//     - equal → found.
//     - greater → col-- (move left).
//     - less → row++ (move down).
//  3. If we walk off the grid, target is absent.
//
// Time:  O(m + n) — at most m+n steps before exiting the grid.
// Space: O(1).
func staircaseSearch(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}
	row := 0                  // start at the top row
	col := len(matrix[0]) - 1 // ...and the rightmost column
	for row < len(matrix) && col >= 0 {
		switch {
		case matrix[row][col] == target:
			return true // exact hit
		case matrix[row][col] > target:
			col-- // current is too big; drop this column, move left
		default:
			row++ // current is too small; drop this row, move down
		}
	}
	return false // walked off the grid without finding target
}

func main() {
	matrix := [][]int{
		{1, 4, 7, 11, 15},
		{2, 5, 8, 12, 19},
		{3, 6, 9, 16, 22},
		{10, 13, 14, 17, 24},
		{18, 21, 23, 26, 30},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(matrix, 5))  // expected true
	fmt.Println(bruteForce(matrix, 20)) // expected false

	fmt.Println("=== Approach 2: Binary Search Each Row ===")
	fmt.Println(binarySearchRows(matrix, 5))  // expected true
	fmt.Println(binarySearchRows(matrix, 20)) // expected false

	fmt.Println("=== Approach 3: Staircase Search from Top-Right (Optimal) ===")
	fmt.Println(staircaseSearch(matrix, 5))  // expected true
	fmt.Println(staircaseSearch(matrix, 20)) // expected false
}
