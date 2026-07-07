package main

import "fmt"

// ── Approach 1: DP 2D Table ───────────────────────────────────────────────────
//
// dpBottomUp solves Unique Paths II (with obstacles) using a 2D DP table.
//
// Intuition:
//   dp[r][c] = number of unique paths to reach (r,c) from (0,0).
//   If obstacleGrid[r][c] == 1: dp[r][c] = 0 (blocked).
//   Otherwise: dp[r][c] = dp[r-1][c] + dp[r][c-1].
//   Edge case: if an obstacle appears on the first row or column, all cells
//   after it in that row/column have 0 paths.
//
// Algorithm:
//   initialize dp[0][0] = 1 (if not blocked)
//   fill first row and column
//   fill rest: dp[r][c] = dp[r-1][c] + dp[r][c-1] if no obstacle else 0
//
// Time:  O(m × n)
// Space: O(m × n)
func dpBottomUp(obstacleGrid [][]int) int {
	m, n := len(obstacleGrid), len(obstacleGrid[0])
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	// starting cell
	if obstacleGrid[0][0] == 1 {
		return 0
	}
	dp[0][0] = 1

	// first column: blocked by any obstacle in the column above
	for r := 1; r < m; r++ {
		if obstacleGrid[r][0] == 0 {
			dp[r][0] = dp[r-1][0]
		}
	}
	// first row
	for c := 1; c < n; c++ {
		if obstacleGrid[0][c] == 0 {
			dp[0][c] = dp[0][c-1]
		}
	}
	// rest of the grid
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			if obstacleGrid[r][c] == 0 {
				dp[r][c] = dp[r-1][c] + dp[r][c-1]
			}
		}
	}
	return dp[m-1][n-1]
}

// ── Approach 2: DP 1D Rolling Array ──────────────────────────────────────────
//
// dpRolling solves Unique Paths II with O(n) space.
//
// Intuition:
//   Same as the rolling-row idea from #62: reuse one array.
//   dp[c] represents paths to the current row's column c.
//   Update in place: dp[c] += dp[c-1] when no obstacle; dp[c] = 0 when obstacle.
//
// Time:  O(m × n)
// Space: O(n)
func dpRolling(obstacleGrid [][]int) int {
	m, n := len(obstacleGrid), len(obstacleGrid[0])
	if obstacleGrid[0][0] == 1 {
		return 0
	}
	dp := make([]int, n)
	dp[0] = 1 // starting cell

	// initialize first row
	for c := 1; c < n; c++ {
		if obstacleGrid[0][c] == 1 {
			dp[c] = 0
		} else {
			dp[c] = dp[c-1]
		}
	}

	for r := 1; r < m; r++ {
		// update first column
		if obstacleGrid[r][0] == 1 {
			dp[0] = 0
		}
		for c := 1; c < n; c++ {
			if obstacleGrid[r][c] == 1 {
				dp[c] = 0
			} else {
				dp[c] += dp[c-1]
			}
		}
	}
	return dp[n-1]
}

func main() {
	fmt.Println("=== Approach 1: DP 2D Table ===")
	g1 := [][]int{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}
	fmt.Printf("grid=%v  got=%d  expected 2\n", g1, dpBottomUp(g1))

	g2 := [][]int{{0, 1}, {0, 0}}
	fmt.Printf("grid=%v  got=%d  expected 1\n", g2, dpBottomUp(g2))

	g3 := [][]int{{1, 0}}
	fmt.Printf("grid=%v  got=%d  expected 0\n", g3, dpBottomUp(g3))

	fmt.Println("=== Approach 2: DP Rolling Array ===")
	g4 := [][]int{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}
	fmt.Printf("grid=%v  got=%d  expected 2\n", g4, dpRolling(g4))

	g5 := [][]int{{0, 1}, {0, 0}}
	fmt.Printf("grid=%v  got=%d  expected 1\n", g5, dpRolling(g5))

	g6 := [][]int{{1, 0}}
	fmt.Printf("grid=%v  got=%d  expected 0\n", g6, dpRolling(g6))
}
