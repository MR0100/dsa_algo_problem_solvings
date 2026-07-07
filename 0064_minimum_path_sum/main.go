package main

import "fmt"

// ── Approach 1: DP 2D Table ───────────────────────────────────────────────────
//
// dpBottomUp solves Minimum Path Sum using a 2D DP table.
//
// Intuition:
//   dp[r][c] = minimum cost to reach cell (r,c) from (0,0).
//   dp[0][0] = grid[0][0]
//   First row: dp[0][c] = dp[0][c-1] + grid[0][c]
//   First col: dp[r][0] = dp[r-1][0] + grid[r][0]
//   Rest: dp[r][c] = min(dp[r-1][c], dp[r][c-1]) + grid[r][c]
//
// Time:  O(m × n)
// Space: O(m × n)
func dpBottomUp(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	dp[0][0] = grid[0][0]
	for r := 1; r < m; r++ {
		dp[r][0] = dp[r-1][0] + grid[r][0]
	}
	for c := 1; c < n; c++ {
		dp[0][c] = dp[0][c-1] + grid[0][c]
	}
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			if dp[r-1][c] < dp[r][c-1] {
				dp[r][c] = dp[r-1][c] + grid[r][c]
			} else {
				dp[r][c] = dp[r][c-1] + grid[r][c]
			}
		}
	}
	return dp[m-1][n-1]
}

// ── Approach 2: DP In-Place (Modify Grid) ────────────────────────────────────
//
// dpInPlace solves Minimum Path Sum by reusing the input grid as the DP table.
//
// Intuition:
//   Same recurrence; write answers back into grid[r][c] directly.
//   Avoids allocating a separate DP table but modifies the input.
//
// Time:  O(m × n)
// Space: O(1)  — no extra allocation (grid reused).
func dpInPlace(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	// first row accumulates left-to-right
	for c := 1; c < n; c++ {
		grid[0][c] += grid[0][c-1]
	}
	// first col accumulates top-to-bottom
	for r := 1; r < m; r++ {
		grid[r][0] += grid[r-1][0]
	}
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			if grid[r-1][c] < grid[r][c-1] {
				grid[r][c] += grid[r-1][c]
			} else {
				grid[r][c] += grid[r][c-1]
			}
		}
	}
	return grid[m-1][n-1]
}

// ── Approach 3: DP Rolling Row ────────────────────────────────────────────────
//
// dpRolling solves Minimum Path Sum with O(n) space.
//
// Intuition:
//   Maintain one row of DP values. To update row r:
//   dp[c] = min(dp[c], dp[c-1]) + grid[r][c]
//   (dp[c] before update = top cell dp[r-1][c]; dp[c-1] after update = left cell dp[r][c-1])
//
// Time:  O(m × n)
// Space: O(n)
func dpRolling(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	dp := make([]int, n)
	dp[0] = grid[0][0]
	for c := 1; c < n; c++ {
		dp[c] = dp[c-1] + grid[0][c] // first row init
	}
	for r := 1; r < m; r++ {
		dp[0] += grid[r][0] // first col: only can come from above
		for c := 1; c < n; c++ {
			top := dp[c] // before update = dp[r-1][c]
			left := dp[c-1] // after update = dp[r][c-1]
			if top < left {
				dp[c] = top + grid[r][c]
			} else {
				dp[c] = left + grid[r][c]
			}
		}
	}
	return dp[n-1]
}

func main() {
	fmt.Println("=== Approach 1: DP 2D Table ===")
	g1 := [][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}
	fmt.Printf("grid=%v  got=%d  expected 7\n", g1, dpBottomUp(g1))
	g2 := [][]int{{1, 2, 3}, {4, 5, 6}}
	fmt.Printf("grid=%v  got=%d  expected 12\n", g2, dpBottomUp(g2))

	fmt.Println("=== Approach 2: DP In-Place ===")
	g3 := [][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}
	fmt.Printf("grid=%v  got=%d  expected 7\n", g3, dpInPlace(g3))
	g4 := [][]int{{1, 2, 3}, {4, 5, 6}}
	fmt.Printf("grid=%v  got=%d  expected 12\n", g4, dpInPlace(g4))

	fmt.Println("=== Approach 3: DP Rolling Row ===")
	g5 := [][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}
	fmt.Printf("grid=%v  got=%d  expected 7\n", g5, dpRolling(g5))
	g6 := [][]int{{1, 2, 3}, {4, 5, 6}}
	fmt.Printf("grid=%v  got=%d  expected 12\n", g6, dpRolling(g6))
}
