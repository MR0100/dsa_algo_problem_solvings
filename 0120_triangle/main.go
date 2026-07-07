package main

import "fmt"

// ── Approach 1: Top-Down DP with Memoization ─────────────────────────────────
//
// minimumTotal solves Triangle: Minimum Path Sum top-down with memoization.
//
// Intuition:
//   At each cell (i,j), the minimum path cost is triangle[i][j] + min of the
//   two cells below it. Memoize to avoid recomputation.
//
// Time:  O(n^2)
// Space: O(n^2) — memo table.
func minimumTotal(triangle [][]int) int {
	n := len(triangle)
	memo := make([][]int, n)
	for i := range memo {
		memo[i] = make([]int, len(triangle[i]))
		for j := range memo[i] {
			memo[i][j] = -1 << 30 // unvisited sentinel
		}
	}

	var dp func(row, col int) int
	dp = func(row, col int) int {
		if row == n-1 {
			return triangle[row][col]
		}
		if memo[row][col] != -1<<30 {
			return memo[row][col]
		}
		left := dp(row+1, col)
		right := dp(row+1, col+1)
		val := triangle[row][col]
		if left < right {
			val += left
		} else {
			val += right
		}
		memo[row][col] = val
		return val
	}
	return dp(0, 0)
}

// ── Approach 2: Bottom-Up DP In-Place (Optimal) ──────────────────────────────
//
// minimumTotalBottomUp solves Triangle: Minimum Path Sum bottom-up.
//
// Intuition:
//   Start from the second-to-last row and work upward.
//   dp[j] = triangle[row][j] + min(dp[j], dp[j+1]).
//   After processing all rows, dp[0] holds the answer.
//
// Time:  O(n^2)
// Space: O(n) — only one row at a time.
func minimumTotalBottomUp(triangle [][]int) int {
	n := len(triangle)
	// copy the last row as starting dp values
	dp := make([]int, n)
	copy(dp, triangle[n-1])

	// build from second-to-last row up
	for row := n - 2; row >= 0; row-- {
		for col := 0; col <= row; col++ {
			if dp[col] < dp[col+1] {
				dp[col] = triangle[row][col] + dp[col]
			} else {
				dp[col] = triangle[row][col] + dp[col+1]
			}
		}
	}
	return dp[0]
}

func main() {
	tri1 := [][]int{{2}, {3, 4}, {6, 5, 7}, {4, 1, 8, 3}}
	tri2 := [][]int{{-10}}

	fmt.Println("=== Approach 1: Top-Down DP ===")
	fmt.Printf("triangle=[[2],[3,4],[6,5,7],[4,1,8,3]]  got=%d  expected 11\n", minimumTotal(tri1))
	fmt.Printf("triangle=[[-10]]  got=%d  expected -10\n", minimumTotal(tri2))

	tri3 := [][]int{{2}, {3, 4}, {6, 5, 7}, {4, 1, 8, 3}}
	tri4 := [][]int{{-10}}

	fmt.Println("=== Approach 2: Bottom-Up DP (Optimal) ===")
	fmt.Printf("triangle=[[2],[3,4],[6,5,7],[4,1,8,3]]  got=%d  expected 11\n", minimumTotalBottomUp(tri3))
	fmt.Printf("triangle=[[-10]]  got=%d  expected -10\n", minimumTotalBottomUp(tri4))
}
