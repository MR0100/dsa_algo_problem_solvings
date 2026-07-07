package main

import "fmt"

// ── Approach 1: Recursion with Memoization ────────────────────────────────────
//
// memoization solves Unique Paths using top-down DP.
//
// Intuition:
//   From each cell (r,c) the robot can move right or down.
//   paths(r,c) = paths(r+1,c) + paths(r,c+1)
//   Base case: any cell on the bottom row or right col = 1 path (only one direction).
//
// Time:  O(m × n) — each cell computed once.
// Space: O(m × n) — memo table + recursion stack.
func memoization(m, n int) int {
	memo := make([][]int, m)
	for i := range memo {
		memo[i] = make([]int, n)
	}
	var dp func(r, c int) int
	dp = func(r, c int) int {
		if r == m-1 || c == n-1 {
			return 1 // bottom row or right column: only one way forward
		}
		if memo[r][c] != 0 {
			return memo[r][c]
		}
		memo[r][c] = dp(r+1, c) + dp(r, c+1)
		return memo[r][c]
	}
	return dp(0, 0)
}

// ── Approach 2: DP Bottom-Up (2D Table) ──────────────────────────────────────
//
// dpBottomUp solves Unique Paths using a 2D DP table.
//
// Intuition:
//   dp[r][c] = number of paths to reach cell (r,c) from (0,0).
//   dp[0][*] = 1, dp[*][0] = 1 (only one path along edges).
//   dp[r][c] = dp[r-1][c] + dp[r][c-1]
//
// Time:  O(m × n)
// Space: O(m × n)
func dpBottomUp(m, n int) int {
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
		dp[i][0] = 1 // left column: only one path (all downs)
	}
	for j := 0; j < n; j++ {
		dp[0][j] = 1 // top row: only one path (all rights)
	}
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			dp[r][c] = dp[r-1][c] + dp[r][c-1]
		}
	}
	return dp[m-1][n-1]
}

// ── Approach 3: DP 1D Rolling Row ────────────────────────────────────────────
//
// dpRolling solves Unique Paths with O(n) space by reusing a single row array.
//
// Intuition:
//   dp[c] starts as the top-row values (all 1s).
//   Each subsequent row update: dp[c] += dp[c-1]
//   (dp[c] = paths from above; dp[c-1] = paths from left after update).
//
// Time:  O(m × n)
// Space: O(n)
func dpRolling(m, n int) int {
	dp := make([]int, n)
	for j := range dp {
		dp[j] = 1 // top row: all 1s
	}
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			dp[c] += dp[c-1] // dp[c] was "from above"; dp[c-1] is "from left"
		}
	}
	return dp[n-1]
}

// ── Approach 4: Combinatorics (Optimal) ──────────────────────────────────────
//
// combinatorics solves Unique Paths using the closed-form formula.
//
// Intuition:
//   The robot always makes exactly (m-1) down moves and (n-1) right moves,
//   total (m+n-2) moves. The number of ways to arrange them is C(m+n-2, m-1).
//
//   C(m+n-2, m-1) = (m+n-2)! / ((m-1)! * (n-1)!)
//
// Time:  O(min(m,n)) — multiply and divide min(m-1,n-1) terms.
// Space: O(1)
func combinatorics(m, n int) int {
	// C(m+n-2, k) where k = min(m-1, n-1) — use the smaller k for efficiency
	total := m + n - 2
	k := m - 1
	if n-1 < k {
		k = n - 1
	}
	result := 1
	for i := 0; i < k; i++ {
		// multiply (total-i) / (i+1) iteratively to avoid overflow
		result = result * (total - i) / (i + 1)
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Memoization ===")
	fmt.Printf("m=3 n=7  got=%d  expected 28\n", memoization(3, 7))
	fmt.Printf("m=3 n=2  got=%d  expected 3\n", memoization(3, 2))
	fmt.Printf("m=1 n=1  got=%d  expected 1\n", memoization(1, 1))

	fmt.Println("=== Approach 2: DP 2D Table ===")
	fmt.Printf("m=3 n=7  got=%d  expected 28\n", dpBottomUp(3, 7))
	fmt.Printf("m=3 n=2  got=%d  expected 3\n", dpBottomUp(3, 2))
	fmt.Printf("m=1 n=1  got=%d  expected 1\n", dpBottomUp(1, 1))

	fmt.Println("=== Approach 3: DP Rolling Row ===")
	fmt.Printf("m=3 n=7  got=%d  expected 28\n", dpRolling(3, 7))
	fmt.Printf("m=3 n=2  got=%d  expected 3\n", dpRolling(3, 2))
	fmt.Printf("m=1 n=1  got=%d  expected 1\n", dpRolling(1, 1))

	fmt.Println("=== Approach 4: Combinatorics ===")
	fmt.Printf("m=3 n=7  got=%d  expected 28\n", combinatorics(3, 7))
	fmt.Printf("m=3 n=2  got=%d  expected 3\n", combinatorics(3, 2))
	fmt.Printf("m=1 n=1  got=%d  expected 1\n", combinatorics(1, 1))
}
