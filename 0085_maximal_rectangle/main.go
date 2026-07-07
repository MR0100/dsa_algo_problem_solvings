package main

import "fmt"

// largestRectangleInHistogram computes the largest rectangle in a histogram.
// (Reused from #84 — monotonic stack approach.)
func largestRectangleInHistogram(heights []int) int {
	heights = append(heights, 0) // sentinel
	stack := []int{}
	maxArea := 0
	for i, h := range heights {
		for len(stack) > 0 && heights[stack[len(stack)-1]] > h {
			topIdx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			height := heights[topIdx]
			var width int
			if len(stack) == 0 {
				width = i
			} else {
				width = i - stack[len(stack)-1] - 1
			}
			if area := height * width; area > maxArea {
				maxArea = area
			}
		}
		stack = append(stack, i)
	}
	return maxArea
}

// ── Approach 1: Histogram per Row (using #84 monotonic stack) ─────────────────
//
// maximalRectangle solves Maximal Rectangle by building a height histogram
// for each row and finding the largest rectangle in that histogram.
//
// Intuition:
//   Treat each row as the base of a histogram. The height of column c in row r
//   is the number of consecutive 1s ending at (r, c) looking upward. If
//   matrix[r][c] == '0', height[c] = 0; otherwise height[c]++.
//
//   Apply the O(n) monotonic stack algorithm from #84 to each row's histogram.
//
// Algorithm:
//   heights = [0] * n
//   for each row r:
//     for each col c:
//       if matrix[r][c] == '1': heights[c]++
//       else: heights[c] = 0
//     maxArea = max(maxArea, largestRectangleInHistogram(heights))
//
// Time:  O(m × n) — m rows, each row O(n) histogram update + O(n) stack pass.
// Space: O(n) — heights array + stack.
func maximalRectangle(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])
	heights := make([]int, n)
	maxArea := 0

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if matrix[r][c] == '1' {
				heights[c]++ // extend height upward
			} else {
				heights[c] = 0 // reset on '0'
			}
		}
		// make a copy since largestRectangleInHistogram modifies by appending
		hCopy := make([]int, n)
		copy(hCopy, heights)
		area := largestRectangleInHistogram(hCopy)
		if area > maxArea {
			maxArea = area
		}
	}
	return maxArea
}

// ── Approach 2: DP (Left, Right, Height arrays) ───────────────────────────────
//
// maximalRectangleDP solves Maximal Rectangle using three DP arrays that
// track, for each cell, the height, leftmost bound, and rightmost bound of
// the maximal rectangle ending at that cell.
//
// Intuition:
//   For each cell (r, c) with matrix[r][c]=='1':
//   - height[c]: consecutive 1s above (including current row).
//   - left[c]: leftmost column index such that all cells from left[c] to c
//     in the current and previous rows (within the current run) are 1.
//   - right[c]: rightmost column index (exclusive) such that all cells from
//     c to right[c]-1 are 1.
//   Area = height[c] * (right[c] - left[c]).
//
// Update rules:
//   height[c] = height[c]+1 if '1' else 0
//   left[c] = max(left[c], curLeft)   where curLeft tracks the start of current '1'-run
//   right[c] = min(right[c], curRight) where curRight tracks the end of current '1'-run
//
// Time:  O(m × n)
// Space: O(n)
func maximalRectangleDP(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])
	height := make([]int, n)
	left := make([]int, n)   // left boundary (inclusive)
	right := make([]int, n)  // right boundary (exclusive)
	for c := range right {
		right[c] = n // initialise to rightmost+1
	}
	maxArea := 0

	for r := 0; r < m; r++ {
		curLeft := 0
		curRight := n

		// update height
		for c := 0; c < n; c++ {
			if matrix[r][c] == '1' {
				height[c]++
			} else {
				height[c] = 0
			}
		}
		// update left boundary (left to right)
		for c := 0; c < n; c++ {
			if matrix[r][c] == '1' {
				if left[c] < curLeft {
					left[c] = curLeft // push left boundary right
				}
			} else {
				left[c] = 0
				curLeft = c + 1 // next run of 1s starts after c
			}
		}
		// update right boundary (right to left)
		for c := n - 1; c >= 0; c-- {
			if matrix[r][c] == '1' {
				if right[c] > curRight {
					right[c] = curRight // push right boundary left
				}
			} else {
				right[c] = n
				curRight = c // next run of 1s ends before c
			}
		}
		// compute area for each column
		for c := 0; c < n; c++ {
			area := height[c] * (right[c] - left[c])
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}

func main() {
	m1 := [][]byte{
		{'1', '0', '1', '0', '0'},
		{'1', '0', '1', '1', '1'},
		{'1', '1', '1', '1', '1'},
		{'1', '0', '0', '1', '0'},
	}
	m2 := [][]byte{{'0'}}
	m3 := [][]byte{{'1'}}

	fmt.Println("=== Approach 1: Histogram per Row ===")
	fmt.Printf("matrix (4×5)  got=%d  expected 6\n", maximalRectangle(m1))

	m1b := [][]byte{
		{'1', '0', '1', '0', '0'},
		{'1', '0', '1', '1', '1'},
		{'1', '1', '1', '1', '1'},
		{'1', '0', '0', '1', '0'},
	}
	fmt.Printf("matrix [[0]]  got=%d  expected 0\n", maximalRectangle(m2))
	fmt.Printf("matrix [[1]]  got=%d  expected 1\n", maximalRectangle(m3))

	fmt.Println("=== Approach 2: DP (Left/Right/Height) ===")
	fmt.Printf("matrix (4×5)  got=%d  expected 6\n", maximalRectangleDP(m1b))
	m2b := [][]byte{{'0'}}
	m3b := [][]byte{{'1'}}
	fmt.Printf("matrix [[0]]  got=%d  expected 0\n", maximalRectangleDP(m2b))
	fmt.Printf("matrix [[1]]  got=%d  expected 1\n", maximalRectangleDP(m3b))
}
