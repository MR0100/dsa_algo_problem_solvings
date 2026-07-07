package main

import "fmt"

// LeetCode 302 — Smallest Rectangle Enclosing Black Pixels
//
// An image is a binary matrix: '0' is a white pixel, '1' is a black pixel. All
// black pixels are connected (a single black region). Given the location (x, y)
// of ONE black pixel, return the area of the smallest axis-aligned rectangle
// that encloses ALL black pixels. Required: sub-O(m·n) runtime.

// ── Approach 1: Brute-Force Scan ─────────────────────────────────────────────
//
// bruteForce solves Smallest Rectangle by scanning every cell and tracking the
// min/max row and column that contain a black pixel.
//
// Intuition:
//
//	The smallest enclosing rectangle is defined by the extreme black pixels:
//	its area is (maxRow − minRow + 1) · (maxCol − minCol + 1). The most direct
//	way to find those extremes is to look at every pixel. Simple and obviously
//	correct — but O(m·n), which violates the follow-up constraint.
//
// Algorithm:
//  1. Initialise minRow/minCol to +∞, maxRow/maxCol to −∞.
//  2. For every black cell, relax the four extremes.
//  3. Return the rectangle area from the extremes.
//
// Time:  O(m·n) — visits every pixel.
// Space: O(1) — four running extremes.
func bruteForce(image []string, x, y int) int {
	if len(image) == 0 || len(image[0]) == 0 {
		return 0
	}
	minRow, maxRow := len(image), -1    // row bounds (init inverted)
	minCol, maxCol := len(image[0]), -1 // col bounds (init inverted)
	for r := 0; r < len(image); r++ {
		for c := 0; c < len(image[0]); c++ {
			if image[r][c] == '1' { // a black pixel widens the bounds
				if r < minRow {
					minRow = r
				}
				if r > maxRow {
					maxRow = r
				}
				if c < minCol {
					minCol = c
				}
				if c > maxCol {
					maxCol = c
				}
			}
		}
	}
	if maxRow == -1 { // no black pixel at all
		return 0
	}
	return (maxRow - minRow + 1) * (maxCol - minCol + 1)
}

// ── Approach 2: Binary Search on Boundaries (Optimal) ────────────────────────
//
// binarySearch solves Smallest Rectangle by binary-searching each of the four
// boundaries, exploiting that black pixels form one connected region.
//
// Intuition:
//
//	Project the black region onto the row axis: the set of rows that contain
//	at least one black pixel is a CONTIGUOUS interval (connectivity guarantees
//	no gaps). The predicate "row r contains a black pixel" is therefore
//	monotone around the region — false above the top edge, true inside, false
//	below the bottom edge. That monotonicity is exactly what binary search
//	needs. Starting from the known black pixel (x, y), we binary-search the
//	top edge in rows [0, x], the bottom edge in rows [x, m], and likewise the
//	left/right column edges. Each search scans one row or column (O(m) or
//	O(n)) per step, giving O(m·log n + n·log m).
//
// Algorithm:
//  1. hasBlackInRow / hasBlackInCol test a whole row/column for any '1'.
//  2. Binary-search the smallest top row that has a black pixel in [0, x].
//  3. Binary-search the first empty row below, i.e. bottom+1 in [x+1, m].
//  4. Do the same for columns using [0, y] and [y+1, n].
//  5. Area = (bottom − top) · (right − left) from the found half-open bounds.
//
// Time:  O(m·log n + n·log m) — each of the four searches costs O(dim·log(other)).
// Space: O(1) — only indices.
func binarySearch(image []string, x, y int) int {
	if len(image) == 0 || len(image[0]) == 0 {
		return 0
	}
	m, n := len(image), len(image[0])

	hasBlackInRow := func(r int) bool { // any '1' in row r?
		for c := 0; c < n; c++ {
			if image[r][c] == '1' {
				return true
			}
		}
		return false
	}
	hasBlackInCol := func(c int) bool { // any '1' in column c?
		for r := 0; r < m; r++ {
			if image[r][c] == '1' {
				return true
			}
		}
		return false
	}

	// searchRows finds a boundary row in [lo, hi). If findFirstBlack is true it
	// returns the smallest row that HAS a black pixel; otherwise the smallest
	// row that has NONE (used to locate the first empty row past the region).
	searchRows := func(lo, hi int, findFirstBlack bool) int {
		for lo < hi {
			mid := (lo + hi) / 2
			if hasBlackInRow(mid) == findFirstBlack {
				hi = mid // condition met → boundary is at or above mid
			} else {
				lo = mid + 1 // condition not met → search below
			}
		}
		return lo
	}
	searchCols := func(lo, hi int, findFirstBlack bool) int {
		for lo < hi {
			mid := (lo + hi) / 2
			if hasBlackInCol(mid) == findFirstBlack {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		return lo
	}

	top := searchRows(0, x, true)       // first black row in [0, x)
	bottom := searchRows(x+1, m, false) // first NON-black row in [x+1, m)
	left := searchCols(0, y, true)      // first black col in [0, y)
	right := searchCols(y+1, n, false)  // first NON-black col in [y+1, n)

	// top/left are inclusive; bottom/right are exclusive (one past the edge).
	return (bottom - top) * (right - left)
}

func main() {
	image := []string{
		"0010",
		"0110",
		"0100",
	}

	fmt.Println("=== Approach 1: Brute-Force Scan ===")
	fmt.Println(bruteForce(image, 0, 2)) // expected 6

	fmt.Println("=== Approach 2: Binary Search on Boundaries (Optimal) ===")
	fmt.Println(binarySearch(image, 0, 2)) // expected 6

	// Single-pixel image: area 1.
	fmt.Println("=== Edge: single black pixel ===")
	fmt.Println(bruteForce([]string{"1"}, 0, 0))   // expected 1
	fmt.Println(binarySearch([]string{"1"}, 0, 0)) // expected 1
}
