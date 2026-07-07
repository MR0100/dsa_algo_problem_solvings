package main

import "fmt"

// ── Approach 1: 2D Bottom-Up DP ──────────────────────────────────────────────
//
// numDistinct solves Distinct Subsequences using 2D DP.
//
// Intuition:
//   dp[i][j] = number of distinct subsequences of s[0:i] that equal t[0:j].
//
//   Recurrence:
//     dp[i][j] = dp[i-1][j]              (don't use s[i-1])
//               + dp[i-1][j-1]  if s[i-1]==t[j-1]  (use s[i-1])
//
//   Base cases:
//     dp[i][0] = 1 for all i (empty t matches any prefix of s once — delete all)
//     dp[0][j] = 0 for j>0  (can't match non-empty t from empty s)
//
// Time:  O(m*n)
// Space: O(m*n)
func numDistinct(s string, t string) int {
	m, n := len(s), len(t)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = 1 // empty t matched by any prefix of s
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			dp[i][j] = dp[i-1][j] // skip s[i-1]
			if s[i-1] == t[j-1] {
				dp[i][j] += dp[i-1][j-1] // use s[i-1]
			}
		}
	}
	return dp[m][n]
}

// ── Approach 2: 1D Rolling DP (Space Optimized) ──────────────────────────────
//
// numDistinctOptimal solves Distinct Subsequences with O(n) space.
//
// Intuition:
//   dp[i][j] depends only on dp[i-1][*]. Roll to a 1D array.
//   Must iterate j from right to left to avoid using updated values.
//
// Time:  O(m*n)
// Space: O(n)
func numDistinctOptimal(s string, t string) int {
	n := len(t)
	dp := make([]int, n+1)
	dp[0] = 1 // empty t: 1 way

	for _, sc := range s {
		// traverse from right to left to avoid using updated values
		for j := n; j >= 1; j-- {
			if sc == rune(t[j-1]) {
				dp[j] += dp[j-1]
			}
		}
	}
	return dp[n]
}

func main() {
	fmt.Println("=== Approach 1: 2D DP ===")
	fmt.Printf("s=rabbbit t=rabbit  got=%d  expected 3\n", numDistinct("rabbbit", "rabbit"))
	fmt.Printf("s=babgbag t=bag  got=%d  expected 5\n", numDistinct("babgbag", "bag"))

	fmt.Println("=== Approach 2: 1D Rolling DP ===")
	fmt.Printf("s=rabbbit t=rabbit  got=%d  expected 3\n", numDistinctOptimal("rabbbit", "rabbit"))
	fmt.Printf("s=babgbag t=bag  got=%d  expected 5\n", numDistinctOptimal("babgbag", "bag"))
}
