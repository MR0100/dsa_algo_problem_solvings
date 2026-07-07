package main

import (
	"fmt"
	"sort"
)

// Best Meeting Point
//
// Given an m x n binary grid where each 1 marks the home of a friend, choose a
// meeting point (any cell) that minimises the total travelling distance, where
// distance between two points is the Manhattan distance
// |p1.x - p2.x| + |p1.y - p2.y|. Return that minimum total distance.

// ── Approach 1: Brute Force (try every cell) ─────────────────────────────────
//
// bruteForce solves Best Meeting Point by testing every grid cell as the
// candidate meeting point and summing Manhattan distances to all homes.
//
// Intuition:
//
//	The meeting point can be any of the m*n cells. For each candidate, compute
//	the total distance to every friend's home and keep the smallest total.
//
// Algorithm:
//  1. Collect all home coordinates.
//  2. For every cell (r, c) in the grid, sum |r-hr| + |c-hc| over all homes.
//  3. Track and return the minimum sum.
//
// Time:  O(m*n*k) — k homes, evaluated at each of the m*n cells.
// Space: O(k) — storing the home coordinates.
func bruteForce(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return 0
	}
	n := len(grid[0])

	type point struct{ r, c int }
	homes := []point{}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] == 1 {
				homes = append(homes, point{r, c}) // record each friend's home
			}
		}
	}

	best := -1 // sentinel: no candidate evaluated yet
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			total := 0
			for _, h := range homes { // sum distance from this cell to every home
				total += abs(r-h.r) + abs(c-h.c)
			}
			if best == -1 || total < best {
				best = total // keep the smallest total distance found
			}
		}
	}
	if best == -1 {
		return 0 // no homes at all
	}
	return best
}

// ── Approach 2: Median via Sorting ───────────────────────────────────────────
//
// medianSort solves Best Meeting Point by exploiting that Manhattan distance
// separates into independent x and y components; the optimal 1-D meeting point
// is the median of the coordinates.
//
// Intuition:
//
//	Total distance = sum|r-hr| + sum|c-hc|. The two sums are independent, so we
//	minimise each separately. In 1-D, the point minimising the sum of absolute
//	distances is the MEDIAN of the points. Collect all row indices and all
//	column indices, sort each, take the median, and sum absolute distances.
//
// Algorithm:
//  1. Gather rows[] and cols[] of every home.
//  2. Sort both arrays.
//  3. distance = minDistance1D(rows) + minDistance1D(cols) where the 1-D helper
//     sums |x - median|.
//
// Time:  O(k log k) — sorting the coordinate lists dominates.
// Space: O(k) — the coordinate lists.
func medianSort(grid [][]int) int {
	rows := []int{}
	cols := []int{}
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			if grid[r][c] == 1 {
				rows = append(rows, r) // collect row of each home
				cols = append(cols, c) // collect column of each home
			}
		}
	}
	sort.Ints(rows) // sort so the middle element is the median
	sort.Ints(cols)
	return minDistance1D(rows) + minDistance1D(cols)
}

// minDistance1D returns the minimum total absolute distance from all points in
// the SORTED slice to their median.
func minDistance1D(sorted []int) int {
	total := 0
	i, j := 0, len(sorted)-1 // walk inward from both ends
	for i < j {
		// The gap between the outer pair must be crossed once, regardless of
		// where inside it the meeting point sits — so add it directly.
		total += sorted[j] - sorted[i]
		i++
		j--
	}
	return total
}

// ── Approach 3: Two-Pointer Without Sorting Columns (Optimal) ────────────────
//
// twoPointerOptimal solves Best Meeting Point by collecting rows in already
// sorted order (row-major scan) and columns in column order, avoiding a sort of
// the rows entirely and pairing extremes with two pointers.
//
// Intuition:
//
//	Scanning the grid row by row yields row indices in non-decreasing order for
//	free. To get columns sorted, scan column by column instead. Then the same
//	two-pointer "sum of outer gaps" trick gives the median distance without an
//	explicit median lookup.
//
// Algorithm:
//  1. Row scan (r outer, c inner): append r for each home -> rows already sorted.
//  2. Column scan (c outer, r inner): append c for each home -> cols already sorted.
//  3. Return twoPointerGap(rows) + twoPointerGap(cols).
//
// Time:  O(m*n) — two full grid scans, no sorting.
// Space: O(k) — the two coordinate lists.
func twoPointerOptimal(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return 0
	}
	n := len(grid[0])

	rows := []int{}
	for r := 0; r < m; r++ { // row-major -> rows come out sorted ascending
		for c := 0; c < n; c++ {
			if grid[r][c] == 1 {
				rows = append(rows, r)
			}
		}
	}

	cols := []int{}
	for c := 0; c < n; c++ { // column-major -> cols come out sorted ascending
		for r := 0; r < m; r++ {
			if grid[r][c] == 1 {
				cols = append(cols, c)
			}
		}
	}
	return twoPointerGap(rows) + twoPointerGap(cols)
}

// twoPointerGap sums the distances between symmetric outer pairs of a SORTED
// slice, which equals the total absolute distance to the median.
func twoPointerGap(sorted []int) int {
	total := 0
	i, j := 0, len(sorted)-1
	for i < j {
		total += sorted[j] - sorted[i] // outer gap crossed exactly once
		i++
		j--
	}
	return total
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	// Example 1:
	//   grid = [[1,0,0,0,1],
	//           [0,0,0,0,0],
	//           [0,0,1,0,0]]  -> 6
	ex1 := [][]int{
		{1, 0, 0, 0, 1},
		{0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0},
	}
	// Example 2:
	//   grid = [[1,1]] -> 1
	ex2 := [][]int{
		{1, 1},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(ex1)) // expected 6
	fmt.Println(bruteForce(ex2)) // expected 1

	fmt.Println("=== Approach 2: Median via Sorting ===")
	fmt.Println(medianSort(ex1)) // expected 6
	fmt.Println(medianSort(ex2)) // expected 1

	fmt.Println("=== Approach 3: Two-Pointer (Optimal) ===")
	fmt.Println(twoPointerOptimal(ex1)) // expected 6
	fmt.Println(twoPointerOptimal(ex2)) // expected 1
}
